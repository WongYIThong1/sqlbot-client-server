package main

import (
	"fmt"
	"sqlbots-client/ui"
)

func main() {
	// 测试界面显示
	fmt.Println("Testing UI components...")
	fmt.Println()
	
	// 测试横幅
	ui.ShowBanner("v1.0")
	
	// 测试登录提示
	ui.ShowLoginPrompt()
	fmt.Println("(simulated input)")
	fmt.Println()
	
	// 测试成功界面
	ui.ShowLoggedIn("testuser", "v1.0")
	
	// 测试错误信息
	ui.ShowError("This is a test error message")
	
	// 测试成功信息
	ui.ShowSuccess("This is a test success message")
}


