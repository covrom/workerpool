package cases

func init() {
	Cases["fastgo"] = Case{
		Prepare: goPrepare,
		Run:     fastgoRun,
	}
}

func fastgoRun(_ int, _ int, amount int, _ chan []byte, chout chan testType) {
	go func() {
		for i := 0; i < amount; i++ {
			go goroutine(processingFast, copyBytes(testData), chout)
		}
	}()

	waitChout(amount, chout)
}
