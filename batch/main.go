package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/terassyi/bpf-batch-test/pkg/bpf"
	"github.com/terassyi/bpf-batch-test/pkg/stats"
)

var memFlag bool

func init() {
	flag.BoolVar(&memFlag, "mem", false, "Show memory stats")
	flag.Parse()
}

func main() {

	m, err := bpf.GetMap()
	if err != nil {
		panic(err)
	}

	fmt.Println("Start batch look up bpf map")

	if memFlag {
		runtime.GC()
		stats.PrintStats("init")
	}

	batchCount := 0
	count := 0

	max := m.MaxEntries()
	const chunk uint32 = 4096
	chSize := int(max / chunk)

	kout := make([]uint32, chunk)
	vout := make([]uint32, chunk)
	var k uint32
	var prev uint32

	for i := 0; i < chSize; i++ {
		c, err := m.BatchLookup(prev, &k, kout, vout, nil)
		if err != nil {
			panic(err)
		}
		count += c
		batchCount++
	}

	if memFlag {
		stats.PrintStats("last")
	}

	fmt.Println("------------------------ ")
	fmt.Println("Finish batch look up bpf map")
	fmt.Printf("Count %d entries\n", count)
	fmt.Printf("Batch Count is %d\n", batchCount)

}
