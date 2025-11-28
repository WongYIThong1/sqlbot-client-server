package ui

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// HideInput 隐藏输入（用于密码/API Key）
func HideInput() (string, error) {
	// Windows 系统
	if runtime.GOOS == "windows" {
		return hideInputWindows()
	}

	// Unix/Linux/Mac 系统
	return hideInputUnix()
}

// hideInputWindows Windows 系统隐藏输入
func hideInputWindows() (string, error) {

	// 获取标准输入句柄
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getStdHandle := kernel32.NewProc("GetStdHandle")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	// STD_INPUT_HANDLE = -10
	handle, _, _ := getStdHandle.Call(uintptr(0xFFFFFFF6))

	// 获取当前控制台模式
	var mode uint32
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))

	// 禁用回显 (ENABLE_ECHO_INPUT = 0x0004)
	var newMode uint32 = mode &^ 0x0004
	setConsoleMode.Call(handle, uintptr(newMode))

	// 读取输入
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 恢复控制台模式
	setConsoleMode.Call(handle, uintptr(mode))

	// 移除换行符
	input = strings.TrimSpace(input)

	return input, nil
}

// hideInputUnix Unix/Linux/Mac 系统隐藏输入
func hideInputUnix() (string, error) {
	// 使用 stty 命令隐藏输入
	exec.Command("stty", "-echo").Run()
	defer exec.Command("stty", "echo").Run()

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 移除换行符
	input = strings.TrimSpace(input)

	return input, nil
}
