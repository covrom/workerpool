package cases

func init() {
	Cases["fastpool"] = Case{
		Prepare: fastpoolPrepare,
		Run:     poolRun,
	}
}

func fastpoolPrepare(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	chin, chout = make(chan []byte, chanLen), make(chan testType)

	for i := 0; i < workers; i++ {
		go worker(processingFast, chin, chout)
	}

	return chin, chout
}
