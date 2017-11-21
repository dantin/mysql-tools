package drc

import (
	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/mysql"
)

// GTIDSet wraps mysql.MysqlGTIDSet.
type GTIDSet struct {
	*mysql.MysqlGTIDSet
}

func parseGTIDSet(gtidStr string) (GTIDSet, error) {
	gs, err := mysql.ParseMysqlGTIDSet(gtidStr)
	if err != nil {
		return GTIDSet{}, errors.Trace(err)
	}
	return GTIDSet{gs.(*mysql.MysqlGTIDSet)}, nil
}

func (g GTIDSet) delete(uuid string) {
	delete(g.Sets, uuid)
}

func (g GTIDSet) contain(uuid string) bool {
	_, ok := g.Sets[uuid]
	return ok
}

func (g GTIDSet) get(uuid string) *mysql.UUIDSet {
	return g.Sets[uuid]
}

func (g GTIDSet) all() map[string]*mysql.UUIDSet {
	return g.Sets
}

func (g GTIDSet) len() int {
	return len(g.Sets)
}

func (g GTIDSet) String() string {
	if g.MysqlGTIDSet == nil {
		return "<nil>"
	}
	return g.MysqlGTIDSet.String()
}
