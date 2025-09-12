package health

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	supabaseClient "github.com/supabase-community/supabase-go"
)

// HealthStatus represents the overall health of the system
type HealthStatus struct {
	Status         string         `json:"status"`
	Timestamp      time.Time      `json:"timestamp"`
	DatabaseHealth DatabaseHealth `json:"database"`
	MemoryHealth   MemoryHealth   `json:"memory"`
	CPUHealth      CPUHealth      `json:"cpu"`
	SystemInfo     SystemInfo     `json:"system"`
}

// DatabaseHealth represents the health of the database connection
type DatabaseHealth struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// MemoryHealth represents the memory usage of the system
type MemoryHealth struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

// CPUHealth represents the CPU load of the system
type CPUHealth struct {
	Cores       int     `json:"cores"`
	LoadAverage float64 `json:"loadAverage"`
}

// SystemInfo provides basic system information
type SystemInfo struct {
	GoVersion    string `json:"goVersion"`
	NumCPU       int    `json:"numCPU"`
	NumGoroutine int    `json:"numGoroutine"`
}

// CheckHealth performs a comprehensive health check
func CheckHealth(client *supabaseClient.Client) HealthStatus {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	// Check Database Connection
	status.DatabaseHealth = checkDatabaseConnection(client)
	if !status.DatabaseHealth.Connected {
		status.Status = "degraded"
	}

	// Check Memory Usage
	status.MemoryHealth = checkMemoryUsage()
	if status.MemoryHealth.UsedPercent > 90 {
		status.Status = "critical"
	}

	// Check CPU Load
	status.CPUHealth = checkCPULoad()
	if status.CPUHealth.LoadAverage > 80 {
		status.Status = "critical"
	}

	// Collect System Info
	status.SystemInfo = getSystemInfo()

	return status
}

// checkDatabaseConnection verifies the database connection
func checkDatabaseConnection(client *supabaseClient.Client) DatabaseHealth {
	dbHealth := DatabaseHealth{
		Connected: false,
	}

	// Attempt to perform a simple query to check connection
	_, _, err := client.
		From("tech_stack").
		Select("id", "", false).
		Limit(1, "").
		Execute()

	if err != nil {
		dbHealth.Error = fmt.Sprintf("Database connection check failed: %v", err)
		return dbHealth
	}

	dbHealth.Connected = true
	return dbHealth
}

// checkMemoryUsage retrieves current memory usage
func checkMemoryUsage() MemoryHealth {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return MemoryHealth{
			Total:       0,
			Used:        0,
			Free:        0,
			UsedPercent: 0,
		}
	}

	return MemoryHealth{
		Total:       vmStat.Total,
		Used:        vmStat.Used,
		Free:        vmStat.Free,
		UsedPercent: vmStat.UsedPercent,
	}
}

// checkCPULoad retrieves current CPU load
func checkCPULoad() CPUHealth {
	cores, err := cpu.Counts(true)
	if err != nil {
		cores = runtime.NumCPU()
	}

	loadAvg, err := cpu.Percent(time.Second, false)
	loadPercent := 0.0
	if err == nil && len(loadAvg) > 0 {
		loadPercent = loadAvg[0]
	}

	return CPUHealth{
		Cores:       cores,
		LoadAverage: loadPercent,
	}
}

// getSystemInfo collects basic system information
func getSystemInfo() SystemInfo {
	return SystemInfo{
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}
}
