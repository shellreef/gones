// Test synchronization using goroutines

package main

import ("fmt"
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
        //fmt.Printf("+") // run PPU
        cycleChannel <- PPU_MASTER_CYCLES

        cycleCount += 1
        if cycleCount == PPU_CYCLES_PER_FRAME {
            nsPerFrame := time.Nanoseconds() - vblankStartedAt
            fps := 1 / (float(nsPerFrame) / 1e9)
            fmt.Printf("%.2f frames/second (%d ns/frame, %d ns/switch)\n", fps, nsPerFrame, nsPerFrame / PPU_CYCLES_PER_FRAME)

            cycleCount = 0
            vblankStartedAt = time.Nanoseconds()
        }

    }
}

func main() {
    //cpuCycleChannel := make(chan int)
    ppuCycleChannel := make(chan int)

    fmt.Printf("Synchronizing at cpu:ppu = %d:%d\n", CPU_MASTER_CYCLES, PPU_MASTER_CYCLES)

    //go cpu(cpuCycleChannel)
    go ppu(ppuCycleChannel)

    masterCycles := 0
    for {
        // for every CPU cycle...
        //masterCycles += <-cpuCycleChannel
        masterCycles += CPU_MASTER_CYCLES
        //fmt.Printf(".")  // run CPU

        // ...run PPU appropriate number of cycles 
        for masterCycles > PPU_MASTER_CYCLES {
            masterCycles -= <-ppuCycleChannel
        }
    }
}

// runtime.GOMAXPROCS(1), default, is fastest (20 fps)
// runtime.GOMAXPROCS(2) slows to 5 fps
// runtime.GOMAXPROCS(4) slows to ~1 fps
