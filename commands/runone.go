package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/covrom/workerpool/cases"
	"github.com/covrom/workerpool/measures"
)

var ErrCaseUndefined = errors.New("case undefined")

func RunOne(name string, workers int, chanLen int, amount int, profile string, nogc bool) {
	log.Printf("runOne: Name %v, Workers %v, ChanLen %v, Amount %v, Profile %v, NoGC %v", name, workers, chanLen, amount, profile, nogc)

	if nogc {
		debug.SetGCPercent(-1)
	}

	if profile != "" {
		f, err := os.Create(profile)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}

		defer pprof.StopCPUProfile()
	}

	c, ok := cases.Cases[name]
	if !ok {
		panic(fmt.Errorf("%q: %w", name, ErrCaseUndefined))
	}

	chin, chout := c.Prepare(workers, chanLen, amount)

	m1, m2 := &runtime.MemStats{}, &runtime.MemStats{}

	start := time.Now()
	runtime.ReadMemStats(m1)

	c.Run(workers, chanLen, amount, chin, chout)

	SpentMs := time.Now().Sub(start)

	runtime.ReadMemStats(m2)

	err := json.NewEncoder(os.Stdout).Encode(
		measures.Measures{
			Case:            name,
			Workers:         workers,
			ChanLen:         chanLen,
			Amount:          amount,
			SpentMs:         SpentMs,
			AllocBytes:      m2.Alloc - m1.Alloc,
			AllocObjects:    m2.Mallocs - m1.Mallocs,
			AllocBytesTotal: m2.TotalAlloc - m1.TotalAlloc,
			SystemMem:       m2.Sys - m1.Sys,
			NoGC:            nogc,
		},
	)
	if err != nil {
		panic(err)
	}
}
