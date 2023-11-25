package stats

import (
	"fmt"
	"runtime"
)

func PrintStats(title string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("----- %s -----\n", title)
	fmt.Printf("Alloc = %d\n", memStats.Alloc)
	fmt.Printf("HeapAlloc = %d\n", memStats.HeapAlloc)
	fmt.Printf("TotalAlloc = %d\n", memStats.TotalAlloc)
	fmt.Printf("Sys = %d\n", memStats.Sys)
	fmt.Printf("NumGC = %d\n", memStats.NumGC)
}
