package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/cilium/ebpf"
	"github.com/terassyi/bpf-batch-test/pkg/bpf"
)

func main() {

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	fmt.Println("Load test program and map")
	m, err := bpf.LoadAndGetMap()
	if err != nil {
		panic(err)
	}

	fmt.Println("Prepare test data...")

	v := uint32(0xdeadbeef)
	for i := 0; i < bpf.MAX_ENTRIES; i++ {
		k := uint32(i)
		if err := m.Update(&k, &v, ebpf.UpdateAny); err != nil {
			panic(err)
		}
	}

	fmt.Println("Finish preparing test data")

	<-ctrlC
	fmt.Println("Unload everything")
	if err := bpf.UnLoad(); err != nil {
		panic(err)
	}

}
