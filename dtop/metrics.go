package main

import (
	"fmt"
	"html"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type Metric struct{ Value int }
type Process struct {
	Name, CPU, Memory string
	cpuValue          float64
}
type Metrics struct {
	Host           string
	CPU, RAM, Disk Metric
	Down, Up       string
	Processes      []Process
}

func collectMetrics() Metrics {
	m := Metrics{Host: hostname(), CPU: Metric{cpuUsage()}, RAM: Metric{memoryUsage()}, Disk: Metric{diskUsage()}}
	m.Down, m.Up = networkRates()
	m.Processes = processes()
	return m
}

func hostname() string {
	h, _ := os.Hostname()
	if h == "" {
		return "localhost"
	}
	return h
}

func cpuUsage() int {
	v, _ := cpu.Percent(100*time.Millisecond, false)
	if len(v) == 0 {
		return 0
	}
	return int(v[0])
}

func memoryUsage() int {
	v, _ := mem.VirtualMemory()
	return int(v.UsedPercent)
}

func diskUsage() int {
	v, _ := disk.Usage("/")
	return int(v.UsedPercent)
}

var networkState struct {
	sync.Mutex
	rx, tx uint64
	at     time.Time
}

func networkRates() (string, string) {
	io, err := net.IOCounters(true)
	if err != nil || len(io) == 0 {
		return "—", "—"
	}
	var rx, tx uint64
	for _, v := range io {
		if v.Name != "lo" && v.Name != "lo0" {
			rx += v.BytesRecv
			tx += v.BytesSent
		}
	}

	now := time.Now()
	networkState.Lock()
	defer networkState.Unlock()
	if networkState.at.IsZero() || rx < networkState.rx || tx < networkState.tx {
		networkState.rx, networkState.tx, networkState.at = rx, tx, now
		return "0 B/s", "0 B/s"
	}
	seconds := now.Sub(networkState.at).Seconds()
	down := uint64(float64(rx-networkState.rx) / seconds)
	up := uint64(float64(tx-networkState.tx) / seconds)
	networkState.rx, networkState.tx, networkState.at = rx, tx, now
	return rate(down), rate(up)
}

func rate(bytes uint64) string {
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB/s", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.0f MB/s", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.0f KB/s", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B/s", bytes)
	}
}

var processState = struct {
	sync.Mutex
	items map[int32]*process.Process
}{items: map[int32]*process.Process{}}

func processes() []Process {
	all, _ := process.Processes()
	processState.Lock()
	defer processState.Unlock()
	rows := make([]Process, 0, len(all))
	for _, current := range all {
		p := processState.items[current.Pid]
		if p == nil {
			p = current
			processState.items[current.Pid] = p
		}
		name, _ := p.Cmdline()
		if strings.TrimSpace(name) == "" {
			name, _ = p.Name()
		}
		cpuValue, _ := p.Percent(0)
		memory, _ := p.MemoryInfo()
		if name == "" || memory == nil {
			continue
		}
		rows = append(rows, Process{lastN(name, 32), fmt.Sprintf("%.1f%%", cpuValue), memorySize(memory.RSS), cpuValue})
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].cpuValue == rows[j].cpuValue {
			return rows[i].Name < rows[j].Name
		}
		return rows[i].cpuValue > rows[j].cpuValue
	})
	if len(rows) > 6 {
		rows = rows[:6]
	}
	return rows
}

func lastN(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return "…" + string(r[len(r)-n+1:])
}

func memorySize(bytes uint64) string {
	if bytes >= 1<<30 {
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	}
	return fmt.Sprintf("%.0f MB", float64(bytes)/(1<<20))
}

func terminalBar(value int) string {
	value = max(0, min(100, value))
	filled := value * 20 / 100
	return strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)
}

func metricsFragment(m Metrics) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<div id="metrics"><p>CPU &nbsp; %s %d%%</p><p>RAM &nbsp; %s %d%%</p><p>DISK &nbsp;%s %d%%</p><p>Network &nbsp; ↓ %s &nbsp; ↑ %s</p><hr><h2>Processes</h2><table><thead><tr><th align="left">COMMAND</th><th align="right">CPU</th><th align="right">MEMORY</th></tr></thead><tbody>`, terminalBar(m.CPU.Value), m.CPU.Value, terminalBar(m.RAM.Value), m.RAM.Value, terminalBar(m.Disk.Value), m.Disk.Value, m.Down, m.Up)
	for _, p := range m.Processes {
		fmt.Fprintf(&b, `<tr><td>%s</td><td align="right">%s</td><td align="right">%s</td></tr>`, html.EscapeString(p.Name), p.CPU, p.Memory)
	}
	b.WriteString(`</tbody></table></div>`)
	return b.String()
}
