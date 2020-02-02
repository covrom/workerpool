package cases

func init() {
	Cases["go"] = Case{
		Prepare: goPrepare,
		Run:     goRun,
	}
}

func goPrepare(_ int, _ int, amount int) (chin chan []byte, chout chan testType) {
	return nil, make(chan testType)
}

func goRun(_ int, _ int, amount int, _ chan []byte, chout chan testType) {
	go runGoroutines(amount, chout)

	waitChout(amount, chout)
}

func goroutine(f func([]byte) testType, b []byte, chout chan<- testType) {
	chout <- f(b)
}

func runGoroutines(amount int, chout chan<- testType) {
	for i := 0; i < amount; i++ {
		go goroutine(processing, copyBytes(testData), chout)
	}
}
