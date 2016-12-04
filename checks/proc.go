package checks

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

const STAT_PATH = "/proc/stat"
const MEMINFO_PATH = "/proc/meminfo"
const FETCH_INTERVAL = 5

const (
	COMMAND_STOP = iota
	COMMAND_CPU  = iota
	COMMAND_MEM  = iota
)

type CpuFetcher struct {
	request chan Request
}

type Request struct {
	command  int
	response chan Result
}

type Result struct {
	Cpu CpuStats
	Mem MemStats
}

type CpuStats struct {
	Total  float64 `json:"total"`
	User   float64 `json:"user"`
	System float64 `json:"system"`
}

type MemStats struct {
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Available uint64 `json:"available"`
	Buffers   uint64 `json:"buffers"`
	Cached    uint64 `json:"cached"`
}

type cpuDiffs struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
	User    uint64
	System  uint64
}

func NewCpuFetcher() *CpuFetcher {
	request := make(chan Request)

	proc := &CpuFetcher{
		request: request,
	}

	return proc
}

func (c *CpuFetcher) Start() {
	go c.fetcher()
}

func (c *CpuFetcher) Stop() {
	c.sendRequest(COMMAND_STOP)
}

func (c *CpuFetcher) GetCpu() CpuStats {
	return c.sendRequest(COMMAND_CPU).Cpu
}

func (c *CpuFetcher) GetMem() MemStats {
	return c.sendRequest(COMMAND_MEM).Mem
}

func (c *CpuFetcher) sendRequest(command int) Result {
	response := make(chan Result)
	req := Request{
		command:  command,
		response: response,
	}

	c.request <- req
	return <-response
}

func (c *CpuFetcher) fetcher() {
	var previous, current cpuDiffs
	lastFetch := time.Now()

	mem, _ := c.fetchMem()

	current, _ = c.fetchCpu()
	time.Sleep(time.Second * FETCH_INTERVAL)

	for {
		now := time.Now()
		if now.Sub(lastFetch) > time.Second*FETCH_INTERVAL {
			var err error

			previous = current
			current, err = c.fetchCpu()
			if err != nil {
				continue
			}

			mem, err = c.fetchMem()
			if err != nil {
				continue
			}

			lastFetch = now
		}

		select {
		case cmd := <-c.request:
			switch cmd.command {
			case COMMAND_STOP:
				return
			case COMMAND_CPU:
				totalDiff := float64(current.Total - previous.Total)
				idleDiff := float64(current.Idle - previous.Idle)
				userDiff := float64(current.User - previous.User)
				systemDiff := float64(current.System - previous.System)

				total := (totalDiff - idleDiff) / totalDiff
				user := userDiff / totalDiff
				system := systemDiff / totalDiff

				result := Result{
					Cpu: CpuStats{
						Total:  total,
						User:   user,
						System: system,
					},
				}
				cmd.response <- result
			case COMMAND_MEM:
				result := Result{
					Mem: mem,
				}
				cmd.response <- result
			}
		case <-time.After(time.Second * FETCH_INTERVAL):
		}
	}
}

func (c *CpuFetcher) fetchCpu() (cpuDiffs, error) {
	stat, err := linuxproc.ReadStat(STAT_PATH)
	if err != nil {
		return cpuDiffs{}, err
	}

	cpu := stat.CPUStatAll

	idle := cpu.Idle + cpu.IOWait
	nonIdle := cpu.User + cpu.System + cpu.Nice + cpu.IRQ + cpu.SoftIRQ + cpu.Steal

	stats := cpuDiffs{
		Total:   idle + nonIdle,
		Idle:    idle,
		NonIdle: nonIdle,
		User:    cpu.User,
		System:  cpu.System,
	}

	return stats, nil
}

func (c *CpuFetcher) fetchMem() (MemStats, error) {
	mem, err := linuxproc.ReadMemInfo(MEMINFO_PATH)
	if err != nil {
		return MemStats{}, err
	}

	stats := MemStats{
		Total:     mem.MemTotal,
		Free:      mem.MemFree,
		Available: mem.MemAvailable,
		Buffers:   mem.Buffers,
		Cached:    mem.Cached,
	}

	return stats, nil
}
