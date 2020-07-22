package utils

import (
	"log"
	"os"
)

var LOG = log.New(os.Stderr, "Error: ", log.LstdFlags | log.Lshortfile)
