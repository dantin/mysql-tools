package replicator

import (
	"context"
	"database/sql"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/dantin/mysql-tools/pkg/sqlutil"
	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go/sync2"
	log "github.com/sirupsen/logrus"
)

const (
	// binlog event timeout
	eventTimeout = 1 * time.Hour

	maxDMLConnectionTimeout = "1m"
)

// Server is the MySQL binlog sync server.
type Server struct {
	sync.Mutex

	cfg *Config

	meta Meta

	syncer *replication.BinlogSyncer

	sourceDB *sql.DB

	closed sync2.AtomicBool

	ctx    context.Context
	cancel context.CancelFunc

	done chan struct{}
}

// NewServer creates a new server.
func NewServer(cfg *Config) *Server {
	svr := &Server{}
	svr.cfg = cfg

	svr.meta = NewLocalMeta(cfg.Meta)
	svr.closed.Set(false)
	svr.ctx, svr.cancel = context.WithCancel(context.Background())
	svr.done = make(chan struct{})

	return svr
}

// Start starts server.
func (s *Server) Start() (err error) {
	// loads meta
	if err = s.meta.Load(); err != nil {
		return errors.Trace(err)
	}

	// run server
	if err := s.run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

// Close closes server.
func (s *Server) Close() {
	s.Lock()
	defer s.Unlock()

	if s.isClosed() {
		return
	}

	s.cancel()
	<-s.done

	// close db connections
	sqlutil.CloseDBs(s.sourceDB)

	if s.syncer != nil {
		s.syncer.Close()
		s.syncer = nil
	}

	s.closed.Set(true)
}

// isClosed checks whether server is closed or not.
func (s *Server) isClosed() bool {
	return s.closed.Get()
}

// run starts main process.
func (s *Server) run() (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("panic. err: %s, stack: %s", e, debug.Stack())
			err = errors.Errorf("pacni error: %v", e)
		}
		// flush jobs
		close(s.done)
	}()

	// create db connections.
	if err = s.createDBs(); err != nil {
		return errors.Trace(err)
	}

	// check format of the binlog.
	if err = sqlutil.CheckBinlogFormat(s.sourceDB); err != nil {
		return errors.Trace(err)
	}

	//
	cfg := replication.BinlogSyncerConfig{
		ServerID: uint32(s.cfg.ServerID),
		Flavor:   "mysql",
		Host:     s.cfg.Source.Host,
		Port:     uint16(s.cfg.Source.Port),
		User:     s.cfg.Source.Username,
		Password: s.cfg.Source.Password,
	}
	s.syncer = replication.NewBinlogSyncer(cfg)

	// TODO check gtid mode
	streamer, _, err := s.getBinlogStreamer()
	if err != nil {
		return errors.Trace(err)
	}

	pos := s.meta.Pos()

	for {
		ctx, cancel := context.WithTimeout(s.ctx, eventTimeout)
		e, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.Canceled {
			// log pos when exiting.
			log.Infof("ready to quit! [%v]", pos)
			return nil
		} else if err == context.DeadlineExceeded {
			log.Info("deadline exceeded.")
			// TODO: sync meta
			continue
		}

		e.Dump(os.Stdout)
	}
}

// createDBs creates connections to DB.
func (s *Server) createDBs() (err error) {
	if s.sourceDB, err = sqlutil.CreateDB(s.cfg.Source, maxDMLConnectionTimeout); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (s *Server) getBinlogStreamer() (*replication.BinlogStreamer, bool, error) {
	if s.cfg.EnableGTID {
		gs, err := s.meta.GTID()
		if err != nil {
			return nil, false, errors.Trace(err)
		}

		streamer, err := s.syncer.StartSyncGTID(gs)
		if err != nil {
			log.Errorf("start sync in gtid mode error %v", err)
			return s.startSyncByPosition()
		}

		return streamer, true, errors.Trace(err)
	}
	return s.startSyncByPosition()
}

func (s *Server) startSyncByPosition() (*replication.BinlogStreamer, bool, error) {
	streamer, err := s.syncer.StartSync(s.meta.Pos())
	return streamer, false, errors.Trace(err)
}
