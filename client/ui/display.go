package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ClearScreen 清屏
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// ShowBanner 显示 ASCII 艺术字横幅
func ShowBanner(version string) {
	banner := `
░██████╗░██████╗░██╗░░░░░██████╗░░█████╗░████████╗░██████╗
██╔════╝██╔═══██╗██║░░░░░██╔══██╗██╔══██╗╚══██╔══╝██╔════╝
╚█████╗░██║██╗██║██║░░░░░██████╦╝██║░░██║░░░██║░░░╚█████╗░
░╚═══██╗╚██████╔╝██║░░░░░██╔══██╗██║░░██║░░░██║░░░░╚═══██╗
██████╔╝░╚═██╔═╝░███████╗██████╦╝╚█████╔╝░░░██║░░░██████╔╝
╚═════╝░░░░╚═╝░░░╚══════╝╚═════╝░░╚════╝░░░░╚═╝░░░╚═════╝░
                   [%s]
`
	fmt.Printf(banner, version)
	fmt.Println()
}

// ShowLoginPrompt 显示登录提示
func ShowLoginPrompt() {
	fmt.Print("please enter your api key: ")
}

// ShowLoggedIn 显示登录成功界面
func ShowLoggedIn(username, version string) {
	ClearScreen()
	ShowBanner(version)
	fmt.Printf("logged in as : %s\n\n", username)
}

// ShowError 显示错误信息
func ShowError(message string) {
	fmt.Printf("\n❌ Error: %s\n\n", message)
}

// ShowSuccess 显示成功信息
func ShowSuccess(message string) {
	fmt.Printf("\n✅ %s\n\n", message)
}
