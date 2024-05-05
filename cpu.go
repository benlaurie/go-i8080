package i8080

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Inverse speed constants; Speed2MHz per instruction is a speed of 2MHz, etc
const (
	SpeedDebug    = 10 * time.Millisecond
	Speed2Mhz     = 500 * time.Nanosecond
	Speed3_125Mhz = 320 * time.Nanosecond
)

// Register index constants
const (
	B uint8 = iota
	C
	D
	E
	H
	L
	M
	A
)

type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

// CPU implements an emulated Intel 8080 CPU
type CPU struct {
	memory    Memory
	Registers [8]uint8
	Flags     flags

	SP     uint16
	PC     uint16
	Halted bool

	ClockTime time.Duration // nanoseconds per clock tick

	bios
	conio
	diskio
}

func (c *CPU) Read(addr uint16) uint8 {
	return c.memory.Read(addr)
}

func (c *CPU) Write(addr uint16, data uint8) {
	c.memory.Write(addr, data)
}

// New creates a new emulated Intel 8080 CPU
func New(conin io.Reader, conout io.Writer, cpmImage []byte, disks []Disk) (c *CPU) {
	c = &CPU{
		ClockTime: Speed2Mhz,
	}
	c.InitBasic(nil)

	c.initBIOS(cpmImage, disks)
	c.initConsole(conin, conout)

	return
}

func (c *CPU) InitBasic(memory Memory) {
	c.Flags = FlagBit1
	c.memory = memory
	c.Halted = false
	c.PC = 0
	c.SP = 0
	for i := 0; i < 8; i++ {
		c.Registers[i] = 0
	}
}

const tickBudget = 10 * time.Millisecond

// Run runs the CPU and returns how many CPU cycles were executed before a halt
func (c *CPU) Run() (cycles uint64) {
	debug := false
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)
	go func() {
		for range sigChan {
			debug = !debug
		}
	}()

	ticker := time.NewTicker(tickBudget)

	defer func() {
		f := recover()
		if f != nil && f != "hlt" {
			panic(f)
		}
	}()
	defer ticker.Stop()

	var timeUsed time.Duration

	nops := 0

	for {
		<-ticker.C // wait for next tick

		for timeUsed < tickBudget {
			if debug {
				fmt.Printf("%8d %s\r\n", cycles, c.Debug())
			}
			op := c.Read(c.PC)
			if op == 0x00 {
				nops++
			} else {
				nops = 0
			}
			c.PC++
			cyclesThisOp := ops[op](op, c)
			cycles += cyclesThisOp
			timeUsed += time.Duration(cyclesThisOp) * c.ClockTime

			if nops > 10 {
				panic("nop")
			}
		}

		timeUsed -= tickBudget
	}
}

func (c *CPU) Step() uint64 {
	if c.Halted {
		return 0
	}
	op := c.Read(c.PC)
	c.PC++
	return ops[op](op, c)
}
