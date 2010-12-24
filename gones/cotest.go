// Test synchronization using goroutines

package main

import ("fmt"
        "runtime"
        "time")


// Synchronization ratio
// 15:5 = 3:1
// 16:5 = 3.2:1
const CPU_MASTER_CYCLES = 15
const PPU_MASTER_CYCLES = 5

const PPU_CYCLES_PER_FRAME = (341 * 262)

// CPU goroutine
func cpu(cycleChannel chan int) {
    for {
        //fmt.Printf("Tick\n")
        cycleChannel <- CPU_MASTER_CYCLES
    }
}

// PPU goroutine
func ppu(cycleChannel chan int) {
    cycleCount := 0
    vblankStartedAt := time.Nanoseconds()

    for {
        //fmt.Printf("Tock\n")
        cycleChannel <- PPU_MASTER_CYCLES

        cycleCount += 1
        if cycleCount == PPU_CYCLES_PER_FRAME {
            nsPerFrame := time.Nanoseconds() - vblankStartedAt
            fps := 1 / (float(nsPerFrame) / 1e9)
            fmt.Printf("%.2f frames/second (%d ns/frame)\n", fps, nsPerFrame)

            cycleCount = 0
            vblankStartedAt = time.Nanoseconds()
        }

    }
}

func main() {
    cpuCycleChannel := make(chan int)
    ppuCycleChannel := make(chan int)

    fmt.Printf("Synchronizing at cpu:ppu = %d:%d\n", CPU_MASTER_CYCLES, PPU_MASTER_CYCLES)

    go cpu(cpuCycleChannel)
    go ppu(ppuCycleChannel)

    masterCycles := 0
    for {
        // for every CPU cycle...
        masterCycles += <-cpuCycleChannel

        // ...run PPU appropriate number of cycles 
        for masterCycles > PPU_MASTER_CYCLES {
            masterCycles -= <-ppuCycleChannel
        }
    }
}

func init() {
    runtime.GOMAXPROCS(4)
}
