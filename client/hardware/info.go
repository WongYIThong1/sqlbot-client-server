package hardware

import (
	"fmt"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// MachineInfo 机器信息结构体
type MachineInfo struct {
	MachineID   string
	MachineName string
	RAM         int // GB
	Cores       int
}

// GetMachineInfo 获取机器信息
func GetMachineInfo() (*MachineInfo, error) {
	// 获取机器唯一 ID
	machineID, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	ramGB := int(memInfo.Total / (1024 * 1024 * 1024)) // 转换为 GB

	// 获取 CPU 核心数
	cpuCount, err := cpu.Counts(true) // true = physical cores
	if err != nil {
		// 如果获取物理核心数失败，尝试逻辑核心数
		cpuCount, err = cpu.Counts(false)
		if err != nil {
			return nil, fmt.Errorf("failed to get CPU count: %w", err)
		}
	}

	return &MachineInfo{
		MachineID:   machineID,
		MachineName: hostname,
		RAM:         ramGB,
		Cores:       cpuCount,
	}, nil
}

