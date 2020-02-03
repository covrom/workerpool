package measures

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Measures struct {
	Case            string        // "Тип запуска",
	Workers         int           // "Воркеры (шт)",
	ChanLen         int           // "Буфер канала (шт)",
	Amount          int           // "Объекты (шт)",
	SpentMs           time.Duration // "Время работы (сек)",
	AllocBytes      uint64        // "Alloc space (байт)",
	AllocObjects    uint64        // "Alloc objects (шт)",
	AllocBytesTotal uint64        // "Total alloc (байт)",
	SystemMem       uint64        // "System memory (байт)",
}

func (m Measures) Fields() []string {
	return []string{
		m.Case,
		fmt.Sprint(m.Workers),
		fmt.Sprint(m.ChanLen),
		fmt.Sprint(m.Amount),
		strings.ReplaceAll(fmt.Sprint(m.SpentMs.Seconds()), ".", ","),
		fmt.Sprint(m.AllocBytes),
		fmt.Sprint(m.AllocObjects),
		fmt.Sprint(m.AllocBytesTotal),
		fmt.Sprint(m.SystemMem),
	}
}

func Fields() []string {
	return []string{
		"Case",            // "Тип запуска",
		"Workers",         // "Воркеры (шт)",
		"ChanLen",         // "Буфер канала (шт)",
		"Amount",          // "Объекты (шт)",
		"SpentMs",           // "Время работы (мс)",
		"AllocBytes",      // "Alloc space (байт)",
		"AllocObjects",    // "Alloc objects (шт)",
		"AllocBytesTotal", // "Total alloc (байт)",
		"SystemMem",       // "System memory (байт)",
	}
}

func Values() []string {
	return []string{
		"SpentMs",           // "Время работы (мс)",
		"AllocBytes",      // "Alloc space (байт)",
		"AllocObjects",    // "Alloc objects (шт)",
		"AllocBytesTotal", // "Total alloc (байт)",
		"SystemMem",       // "System memory (байт)",
	}
}

var ErrNoSuchField = errors.New("no such field")

func (m Measures) Value(n string) float64 {
	switch n {
	case "SpentMs":
		return float64(m.SpentMs) / float64(time.Millisecond)
	case "AllocBytes":
		return float64(m.AllocBytes)
	case "AllocObjects":
		return float64(m.AllocObjects)
	case "AllocBytesTotal":
		return float64(m.AllocBytesTotal)
	case "SystemMem":
		return float64(m.SystemMem)
	default:
		panic(fmt.Errorf("%q: %w", n, ErrNoSuchField))
	}
}
