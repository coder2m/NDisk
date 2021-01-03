package app

import (
	xapp "github.com/myxy99/component"
	"ndisk/cmd"
	s "ndisk/internal/getway"
)

func Run(stopCh <-chan struct{}) error {
	server := NewServer()
	err := server.PrepareRun(stopCh)
	if err != nil {
		return err
	}
	return server.Run(stopCh)
}

func NewServer() cmd.App {
	xapp.PrintVersion()
	return &s.Server{}
}
