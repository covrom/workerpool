package cases

func init() {
	Cases["pool"] = Case{
		Prepare: poolPrepare,
		Run:     poolRun,
	}
}

func poolPrepare(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	chin, chout = make(chan []byte, chanLen), make(chan testType, chanLen)

	for i := 0; i < workers; i++ {
		go worker(processing, chin, chout)
	}

	return chin, chout
}

func poolRun(workers int, chanLen int, amount int, chin chan []byte, chout chan testType) {
	go func() {
		for i := 0; i < amount; i++ {
			chin <- copyBytes(testData)
		}
	}()

	waitChout(amount, chout)
}

func worker(f func([]byte) testType, chin <-chan []byte, chout chan<- testType) {
	for b := range chin {
		chout <- f(b)
	}
}
