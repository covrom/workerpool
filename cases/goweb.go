package cases

import (
	"bytes"
	"io"
	"net/http"
)

func init() {
	Cases["goweb"] = Case{
		Prepare: goPrepareWeb,
		Run:     goRunWeb,
	}
}

func goPrepareWeb(_ int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewReader(testData))
	})
	go http.ListenAndServe("127.0.0.1:9000", nil)
	return nil, make(chan testType, chanLen)
}

func goRunWeb(_ int, _ int, amount int, _ chan []byte, chout chan testType) {
	go func() {
		for i := 0; i < amount; i++ {
			go goroutineWeb(processingWeb, nil, chout)
		}
	}()

	waitChout(amount, chout)
}

func goroutineWeb(f func() testType, b []byte, chout chan<- testType) {
	chout <- f()
}
