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

	fmt.Println("Start to iterate bpf map")

	if memFlag {
		runtime.GC()
		stats.PrintStats("init")
	}

	var k, v uint32
	count := 0

	iter := m.Iterate()
	for iter.Next(&k, &v) {
		count++
	}

	if memFlag {
		stats.PrintStats("last")
	}

	fmt.Println("------------------------ ")
	fmt.Println("Finish iterating bpf map")
	fmt.Printf("Count %d entries\n", count)

}
