package monitoring

import (
	"launchpad.net/gocheck"
	"testing"
)

func TestWithGocheck(t *testing.T) { gocheck.TestingT(t) }

type TestSuite struct{}

var _ = gocheck.Suite(&TestSuite{})
