// +build !windows

package dashboard_models

import (
	"github.com/shirou/gopsutil/mem"
	"sync"
)

type Memory struct {
	sync.Mutex
	SwapTotal uint64 `json:"swap_total"`
	SwapFree  uint64 `json:"swap_free"`
	MemTotal  uint64 `json:"mem_total"`
	MemFree   uint64 `json:"mem_free"`
}

func (m *Memory) Update() {

	m.Lock()
	swap, _ := mem.SwapMemory()
	m.SwapFree = swap.Free / 1024
	m.SwapTotal = swap.Total / 1024

	memory, _ := mem.VirtualMemory()
	m.MemFree = memory.Free / 1024
	m.MemTotal = memory.Total / 1024
	m.Unlock()
}
