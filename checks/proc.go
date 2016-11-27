package checks

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

const STAT_PATH = "/proc/stat"
const FETCH_INTERVAL = 5

const (
	COMMAND_GET  = iota
	COMMAND_STOP = iota
)

type CpuFetcher struct {
	command  chan int
	response chan Result
}

type Result struct {
	Total  float64 `json:"total"`
	User   float64 `json:"user"`
	System float64 `json:"system"`
}

type cpuStats struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
	User    uint64
	System  uint64
}

func NewCpuFetcher() *CpuFetcher {
	command := make(chan int)
	response := make(chan Result)

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

func (c *CpuFetcher) Get() Result {
	c.command <- COMMAND_GET
	return <-c.response
}

func (c *CpuFetcher) fetcher() {
	var previous, current cpuStats
	lastFetch := time.Now()

	current, _ = c.cpuTime()
	time.Sleep(time.Second * FETCH_INTERVAL)

	for {
		now := time.Now()
		if now.Sub(lastFetch) > time.Second*FETCH_INTERVAL {
			var err error

			previous = current
			current, err = c.cpuTime()
			if err != nil {
				continue
			}
			lastFetch = now
		}

		select {
		case cmd := <-c.command:
			switch cmd {
			case COMMAND_STOP:
				return
			case COMMAND_GET:
				totalDiff := float64(current.Total - previous.Total)
				idleDiff := float64(current.Idle - previous.Idle)
				userDiff := float64(current.User - previous.User)
				systemDiff := float64(current.System - previous.System)

				total := (totalDiff - idleDiff) / totalDiff
				user := userDiff / totalDiff
				system := systemDiff / totalDiff

				result := Result{
					Total:  total,
					User:   user,
					System: system,
				}
				c.response <- result
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
		User:    cpu.User,
		System:  cpu.System,
	}

	return stats, nil
}
