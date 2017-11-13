package logutil

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testLogSuite{})

type testLogSuite struct {
	buf *bytes.Buffer
}

func (s *testLogSuite) SetUpSuite(c *C) {
	s.buf = &bytes.Buffer{}
}

func (s *testLogSuite) TestStringToLogLevel(c *C) {
	c.Assert(stringToLogLevel("fatal"), Equals, log.FatalLevel)
	c.Assert(stringToLogLevel("ERROR"), Equals, log.ErrorLevel)
	c.Assert(stringToLogLevel("warn"), Equals, log.WarnLevel)
	c.Assert(stringToLogLevel("warning"), Equals, log.WarnLevel)
	c.Assert(stringToLogLevel("debug"), Equals, log.DebugLevel)
	c.Assert(stringToLogLevel("info"), Equals, log.InfoLevel)
	c.Assert(stringToLogLevel("else"), Equals, log.InfoLevel)
}
