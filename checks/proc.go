package checks

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

const STAT_PATH = "/proc/stat"
const WARMUP_INTERVAL = 1
const FETCH_INTERVAL = 5

const (
	COMMAND_GET  = iota
	COMMAND_STOP = iota
)

type CpuFetcher struct {
	command  chan int
	response chan float64
}

type cpuStats struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
}

func NewCpuFetcher() *CpuFetcher {
	command := make(chan int)
	response := make(chan float64)

	proc := &CpuFetcher{
		command:  command,
		response: response,
	}

	return proc
}

func (c *CpuFetcher) Start() {
	go c.fetcher()
}

func (c *CpuFetcher) Stop() {
	c.command <- COMMAND_STOP
}

func (c *CpuFetcher) GetTotal() float64 {
	c.command <- COMMAND_GET
	return <-c.response
}

func (c *CpuFetcher) fetcher() {
	var previous cpuStats

	for {
		cpu, err := c.cpuTime()
		if err != nil {
			continue
		}

		if previous.Idle == 0 {
			previous = cpu
			time.Sleep(time.Second * WARMUP_INTERVAL)
			continue
		}

		totalDiff := float64(cpu.Total - previous.Total)
		idleDiff := float64(cpu.Idle - previous.Idle)

		percentage := (totalDiff - idleDiff) / totalDiff

		select {
		case cmd := <-c.command:
			switch cmd {
			case COMMAND_STOP:
				return
			case COMMAND_GET:
				c.response <- percentage
			}
		case <-time.After(time.Second * FETCH_INTERVAL):
		}
	}
}

func (c *CpuFetcher) cpuTime() (cpuStats, error) {
	stat, err := linuxproc.ReadStat(STAT_PATH)
	if err != nil {
		return cpuStats{}, err
	}

	cpu := stat.CPUStatAll

	idle := cpu.Idle + cpu.IOWait
	nonIdle := cpu.User + cpu.System + cpu.Nice + cpu.IRQ + cpu.SoftIRQ + cpu.Steal

	stats := cpuStats{
		Total:   idle + nonIdle,
		Idle:    idle,
		NonIdle: nonIdle,
	}

	return stats, nil
}
