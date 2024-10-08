package main

import "fmt"

var (
	ErrExistingServer = fmt.Errorf("server already started")
	ErrClosedServer   = fmt.Errorf("server already closed")

	ErrUnknownBot = fmt.Errorf("bot does not exist")
)
