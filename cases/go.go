package cases

func init() {
	Cases["go"] = Case{
		Prepare: goPrepare,
		Run:     goRun,
	}
}

func goPrepare(_ int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	return nil, make(chan testType, chanLen)
}

func goRun(_ int, _ int, amount int, _ chan []byte, chout chan testType) {
	go func() {
		for i := 0; i < amount; i++ {
			go goroutine(processing, copyBytes(testData), chout)
		}
	}()

	waitChout(amount, chout)
}

func goroutine(f func([]byte) testType, b []byte, chout chan<- testType) {
	chout <- f(b)
}
