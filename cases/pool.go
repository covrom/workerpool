package cases

import (
	"errors"
	"fmt"
	"log"
)

func init() {
	Cases["pool"] = Case{
		Prepare: workersPrepare,
		Run:     workersRun,
	}
}

func workersPrepare(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	chin, chout = make(chan []byte, chanLen), make(chan testType, chanLen)

	runWorkers(workers, chin, chout)

	return chin, chout
}

func workersRun(workers int, chanLen int, amount int, chin chan []byte, chout chan testType) {
	go useWorkers(amount, chin)

	waitChout(amount, chout)
}

func worker(f func([]byte) testType, chin <-chan []byte, chout chan<- testType) {
	for b := range chin {
		chout <- f(b)
	}
}

func runWorkers(workers int, chin <-chan []byte, chout chan<- testType) {
	log.Printf("starting %d workers", workers)

	for i := 0; i < workers; i++ {
		go worker(processing, chin, chout)
	}

	log.Printf("started %d workers", workers)
}

func useWorkers(amount int, chin chan<- []byte) {
	log.Printf("writing %d inputs", amount)

	for i := 0; i < amount; i++ {
		chin <- copyBytes(testData)
	}

	close(chin)

	log.Printf("writed %d inputs", amount)
}

var (
	ErrUnexpectedResult = errors.New("unexpected result")
	ErrUnreachable      = errors.New("unreachable reached")
)

func waitChout(amount int, chout <-chan testType) {
	log.Printf("reading output")

	for obj := range chout {
		_ = obj
		//if !testResult.Equal(obj) {
		//	panic(fmt.Errorf("got %+v, expect %+v: %w", obj, testResult, ErrUnexpectedResult))
		//}
		amount--
		if amount == 0 {
			log.Printf("completed output")
			return
		}
	}

	panic(fmt.Errorf("%d: %w", amount, ErrUnreachable))
}
