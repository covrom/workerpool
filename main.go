package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

var (
	head         = flag.Bool("head", false, "print csv header")
	pooltyp      = flag.Bool("pool", false, "workerpool type")
	pooltypf     = flag.Bool("fullpool", false, "full filled buffer before workerpool type")
	gotyp        = flag.Bool("go", false, "goroutine type")
	gotypd       = flag.Bool("dgo", false, "delayed goroutine type")
	gotypf       = flag.Bool("fgo", false, "fast goroutine type")
	fastpooltyp  = flag.Bool("fastpool", false, "fast workerpool type")
	numWorkers   = flag.Int("w", 20, "number of workers for workerpool types")
	chanSize     = flag.Int("ch", 200, "channel buffer size for workerpool type")
	totalObjects = flag.Int("c", 30e4, "objects count")
)

type Object struct {
	A uuid.UUID
	B float64
	C string
}

var chin chan []byte
var chnull chan Object

func processing(b []byte) {
	obj := Object{}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	chnull <- obj
}

func worker() {
	for jo := range chin {
		processing(jo)
	}
}

func runWorkers() {
	for i := 0; i < *numWorkers; i++ {
		go worker()
	}
}

func useWorkers() {
	for i := 0; i < *totalObjects; i++ {
		b := make([]byte, len(testCase))
		copy(b, testCase)
		chin <- b
	}
	close(chin)
}

func waitWorkers() {
	i := 0
	for obj := range chnull {
		i++
		if i == *totalObjects {
			if obj.A.String() != "f47ac10b-58cc-0372-8567-0e02b2c3d479" ||
				obj.B != 3.14159265359 ||
				obj.C != "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Dapibus ultrices in iaculis nunc sed. At erat pellentesque adipiscing commodo elit at imperdiet dui accumsan. Dignissim sodales ut eu sem. Mattis vulputate enim nulla aliquet porttitor lacus luctus." {
				log.Fatalf("obj = %+v\n", obj)
			}
			return
		}
	}
}

func simpleGo() {
	for i := 0; i < *totalObjects; i++ {
		b := make([]byte, len(testCase))
		copy(b, testCase)
		go processing(b)
	}
}

func processingDelayed(b []byte, chstart chan bool) {
	<-chstart
	obj := Object{}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	chnull <- obj
}

func delayedGo() {
	chstart := make(chan bool)
	for i := 0; i < *totalObjects; i++ {
		b := make([]byte, len(testCase))
		copy(b, testCase)
		go processingDelayed(b, chstart)
	}
	close(chstart)
}

var poolBuf = &sync.Pool{
	New: func() interface{} { return make([]byte, len(testCase)) },
}

var poolParser fastjson.ParserPool

func processingFast(b []byte) {
	obj := Object{}
	p := poolParser.Get()
	v, err := p.ParseBytes(b)
	if err != nil {
		log.Fatal(err)
	}
	obj.A, err = uuid.ParseBytes(v.GetStringBytes("A"))
	if err != nil {
		log.Fatal(err)
	}
	obj.B = v.GetFloat64("B")
	obj.C = string(v.GetStringBytes("C"))
	poolBuf.Put(b)
	poolParser.Put(p)
	chnull <- obj
}

func fastGo() {
	for i := 0; i < *totalObjects; i++ {
		b := poolBuf.Get().([]byte)
		copy(b, testCase)
		go processingFast(b)
	}
}

func fastWorker() {
	for jo := range chin {
		processingFast(jo)
	}
}

func runFastWorkers() {
	for i := 0; i < *numWorkers; i++ {
		go fastWorker()
	}
}

func main() {
	flag.Parse()

	chin = make(chan []byte, *chanSize)
	chnull = make(chan Object, *chanSize)

	var mode string

	m1 := &runtime.MemStats{}
	m2 := &runtime.MemStats{}
	var start time.Time

	runtime.ReadMemStats(m1)

	switch {
	case *pooltyp:

		mode = "Пул воркеров"

		runWorkers()
		start = time.Now()
		go useWorkers()
		waitWorkers()

	case *pooltypf:

		mode = "Пул воркеров (буфер заполнен)"
		*chanSize = *totalObjects
		chin = make(chan []byte, *totalObjects)

		useWorkers()
		start = time.Now()
		runWorkers()
		waitWorkers()

	case *gotyp:

		mode = "Горутины"
		*numWorkers = 0

		start = time.Now()
		go simpleGo()
		waitWorkers()

	case *gotypd:

		mode = "Горутины (одновременно)"
		*numWorkers = *totalObjects
		delayedGo()
		start = time.Now()
		waitWorkers()

	case *gotypf:

		mode = "Горутины (fast)"
		*numWorkers = 0

		start = time.Now()
		go fastGo()
		waitWorkers()

	case *fastpooltyp:

		mode = "Пул воркеров (fast)"

		runFastWorkers()
		start = time.Now()
		go useWorkers()
		waitWorkers()

	default:
		flag.Usage()
		os.Exit(1)
	}

	runtime.ReadMemStats(m2)

	w := csv.NewWriter(os.Stdout)
	w.Comma = ';'

	if *head {
		w.Write([]string{
			"Тип запуска",
			"Воркеры (шт)",
			"Буфер канала (шт)",
			"Объекты (шт)",
			"Время работы (сек)",
			"Alloc space (байт)",
			"Alloc objects (шт)",
			"Total alloc (байт)",
			"System memory (байт)",
		})
	}

	w.Write([]string{
		mode,
		fmt.Sprint(*numWorkers),
		fmt.Sprint(*chanSize),
		fmt.Sprint(*totalObjects),
		strings.ReplaceAll(fmt.Sprint(float64(time.Since(start))/float64(time.Second)), ".", ","),
		fmt.Sprint(m2.Alloc - m1.Alloc),
		fmt.Sprint(m2.Mallocs - m1.Mallocs),
		fmt.Sprint(m2.TotalAlloc - m1.TotalAlloc),
		fmt.Sprint(m2.Sys - m1.Sys),
	})

	w.Flush()
}

var testCase = []byte(`{
	"A":"f47ac10b-58cc-0372-8567-0e02b2c3d479",
	"C":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Dapibus ultrices in iaculis nunc sed. At erat pellentesque adipiscing commodo elit at imperdiet dui accumsan. Dignissim sodales ut eu sem. Mattis vulputate enim nulla aliquet porttitor lacus luctus.",
	"B":3.14159265359
}`)
