package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", log.Lshortfile)
