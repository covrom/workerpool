package cases

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

var (
	poolBuf    = &sync.Pool{New: func() interface{} { return make([]byte, len(testData)) }}
	poolParser fastjson.ParserPool
)

func copyBytes(src []byte) []byte {
	b := poolBuf.Get().([]byte)
	copy(b, src)
	return b
}

func processing(b []byte) testType {
	defer poolBuf.Put(b)

	obj := testType{}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return obj
}

func processingFast(b []byte) testType {
	defer poolBuf.Put(b)

	p := poolParser.Get()
	defer poolParser.Put(p)

	obj := testType{}

	v, err := p.ParseBytes(b)
	if err != nil {
		panic(err)
	}

	obj.A, err = uuid.ParseBytes(v.GetStringBytes("A"))
	if err != nil {
		panic(err)
	}

	obj.B = v.GetFloat64("B")
	obj.C = string(v.GetStringBytes("C"))

	return obj
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
