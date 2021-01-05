package app

import (
	"github.com/myxy99/ndisk/cmd"
	s "github.com/myxy99/ndisk/internal/nuser"
	"sync"
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
