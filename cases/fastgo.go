package cases

func init() {
	Cases["fastgo"] = Case{
		Prepare: goPrepare,
		Run:     fastgoRun,
	}
}

func fastgoRun(_ int, _ int, amount int, _ chan []byte, chout chan testType) {
	go runFastGo(amount, chout)

	waitChout(amount, chout)
}

func fastgo(f func([]byte) testType, b []byte, chout chan<- testType) {
	chout <- f(b)
}

func runFastGo(amount int, chout chan<- testType) {
	for i := 0; i < amount; i++ {
		go fastgo(processingFast, copyBytes(testData), chout)
	}
}
