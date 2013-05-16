package endpoint

import (
	"log"
	"os"
)

const VERSION = "1"

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
