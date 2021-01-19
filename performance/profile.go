package performance

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func StartCpuProfile() {
	f, err := os.Create("cpuprofile.prof")
	if err != nil {
		fmt.Println("could not create CPU profile: ", err)
		return
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile: ", err)
	}
}

func StopCpuProfile() {
	pprof.StopCPUProfile()
}

func StartMemoryProfile() {
	f, err := os.Create("memprofile.prof")
	if err != nil {
		fmt.Println("could not create memory profile: ", err)
		return
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Println("could not write memory profile: ", err)
	}
	f.Close()
}
