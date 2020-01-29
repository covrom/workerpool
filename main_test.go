package main

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func BenchmarkPool(b *testing.B) { // CPU util: 90%
	*totalObjects = b.N
	chin = make(chan []byte, 200)
	chnull = make(chan Object, 200)
	runWorkers()
	// b.ResetTimer()
	go useWorkers()
	waitWorkers()
}

func BenchmarkPoolFull(b *testing.B) { // CPU util: 100%
	*totalObjects = b.N
	chin = make(chan []byte, *totalObjects)
	chnull = make(chan Object, 200)
	useWorkers()
	// b.ResetTimer()
	runWorkers()
	waitWorkers()
}

func BenchmarkGo(b *testing.B) { // CPU util: 96%
	*totalObjects = b.N
	chin = make(chan []byte, 200)
	chnull = make(chan Object, 200)
	go simpleGo()
	waitWorkers()
}

func BenchmarkGoAll(b *testing.B) { // CPU util: 87%
	*totalObjects = b.N
	chin = make(chan []byte, 200)
	chnull = make(chan Object, 200)
	delayedGo()
	// b.ResetTimer()
	waitWorkers()
}

func BenchmarkGoFast(b *testing.B) { // CPU util: 93%
	*totalObjects = b.N
	chin = make(chan []byte, 200)
	chnull = make(chan Object, 200)
	go fastGo()
	waitWorkers()
}

func BenchmarkPoolFast(b *testing.B) { // CPU util: 90%
	*totalObjects = b.N
	chin = make(chan []byte, 200)
	chnull = make(chan Object, 200)
	runFastWorkers()
	// b.ResetTimer()
	go useWorkers()
	waitWorkers()
}