package cmd

type App interface {
	PrepareRun(stopCh <-chan struct{}) error
	Run(stopCh <-chan struct{}) error
}
