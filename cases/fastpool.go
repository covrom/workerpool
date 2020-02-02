package cases

import (
	"log"
)

func init() {
	Cases["fastpool"] = Case{
		Prepare: fastpoolRun,
		Run:     workersRun,
	}
}

func fastpoolRun(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	chin, chout = make(chan []byte, chanLen), make(chan testType)

	runFastWorkers(workers, chin, chout)

	return chin, chout
}

func runFastWorkers(workers int, chin <-chan []byte, chout chan<- testType) {
	log.Printf("starting %d workers", workers)

	for i := 0; i < workers; i++ {
		go worker(processingFast, chin, chout)
	}

	log.Printf("started %d workers", workers)
}
