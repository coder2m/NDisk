package app

import (
	"sync"

	"github.com/coder2z/ndisk/cmd"
	s "github.com/coder2z/ndisk/internal/authority"
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
