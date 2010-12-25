// Test synchronization using goroutines

package main

import ("fmt"
        "time")


// Synchronization ratio
// See http://nesdev.parodius.com/bbs/viewtopic.php?p=61467
// 15:5 = 12:4 = 3:1
// 16:5 = 3.2:1
const CPU_MASTER_CYCLES = 12
const PPU_MASTER_CYCLES = 4

const PPU_CYCLES_PER_FRAME = (341 * 262)

// 1 / (1.79 MHz * 3)
const PIXELS_NS_TARGET = 187   // ns/pixel

// CPU goroutine
func cpu(cycleChannel chan int) {
    for {
        //fmt.Printf(".")  // run CPU
        cycleChannel <- CPU_MASTER_CYCLES
    }
}

var cycleCount = 0
var vblankStartedAt = time.Nanoseconds()

// Run one cycle of PPU
func ppuTick() {
    //fmt.Printf("+") // run PPU
    cycleCount += 1
    if cycleCount == PPU_CYCLES_PER_FRAME {
        nsPerFrame := time.Nanoseconds() - vblankStartedAt
        fps := 1 / (float(nsPerFrame) / 1e9)
        fmt.Printf("%.2f frames/second (%d ns/frame, %d ns/pixel (%d))\n", fps, nsPerFrame, nsPerFrame / PPU_CYCLES_PER_FRAME,
                PIXELS_NS_TARGET - nsPerFrame / PPU_CYCLES_PER_FRAME)

        cycleCount = 0
        vblankStartedAt = time.Nanoseconds()
    }
}

func main() {
    cpuCycleChannel := make(chan int)

    fmt.Printf("Synchronizing at ppu:cpu = %d:%d\n", CPU_MASTER_CYCLES, PPU_MASTER_CYCLES)

    go cpu(cpuCycleChannel)
    //go ppu(ppuCycleChannel)
    
    masterCycles := 0

    for {
        // for every CPU cycle...
        //fmt.Printf(".")  // run CPU
        //masterCycles := <-cpuCycleChannel
        masterCycles += <-cpuCycleChannel

        // ...run PPU appropriate number of cycles 
        //ppuCycleChannel <- masterCycles / PPU_MASTER_CYCLES

        for masterCycles > PPU_MASTER_CYCLES {
            ppuTick()
            masterCycles -= PPU_MASTER_CYCLES
        }
    }
}

// runtime.GOMAXPROCS(1), default, is fastest (20 fps)
// runtime.GOMAXPROCS(2) slows to 5 fps
// runtime.GOMAXPROCS(4) slows to ~1 fps
