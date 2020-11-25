package utils

import (
	"log"
	"os"
)

// A logger used to log errors into STDERR
// The LOG defined here is imported in other packages and Println method is used to log
var LOG = log.New(os.Stderr, "Error: ", log.LstdFlags|log.Lshortfile)
