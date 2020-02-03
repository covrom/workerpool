package cases

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

func processing(b []byte) testType {
	obj := testType{}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return obj
}

var poolParser fastjson.ParserPool

func processingFast(b []byte) testType {
	obj := testType{}
	p := poolParser.Get()

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
