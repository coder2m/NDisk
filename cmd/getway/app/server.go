package app

import (
	"sync"

	"github.com/coder2m/ndisk/cmd"
	s "github.com/coder2m/ndisk/internal/getway"
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
	return &s.Server{
		WaitGroup: new(sync.WaitGroup),
	}
}
