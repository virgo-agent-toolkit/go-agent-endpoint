package endpoint

import (
	"log"
	"os"
)

// VERSION is the protocol version
const VERSION = "1"

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
