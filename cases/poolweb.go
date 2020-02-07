package cases

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	Cases["poolweb"] = Case{
		Prepare: poolPrepareWeb,
		Run:     poolRunWeb,
	}
}

func poolPrepareWeb(workers int, chanLen int, amount int) (chin chan []byte, chout chan testType) {
	chin, chout = make(chan []byte, chanLen), make(chan testType, chanLen)

	for i := 0; i < workers; i++ {
		go workerWeb(processingWeb, chin, chout)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewReader(testData))
	})
	go http.ListenAndServe("127.0.0.1:9000", nil)

	return chin, chout
}

func poolRunWeb(workers int, chanLen int, amount int, chin chan []byte, chout chan testType) {
	go func() {
		for i := 0; i < amount; i++ {
			chin <- nil
		}
	}()

	waitChout(amount, chout)
}

func workerWeb(f func() testType, chin <-chan []byte, chout chan<- testType) {
	for _ = range chin {
		chout <- f()
	}
}

func processingWeb() testType {
	resp, err := http.Get("http://127.0.0.1:9000/")
	if err != nil {
		log.Printf("processingWeb error: %s", err)
		return testType{}
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("processingWeb error: %s", err)
		return testType{}
	}

	obj := testType{}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return obj
}
