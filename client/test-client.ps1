# 测试客户端脚本
$env:ENCRYPTION_KEY = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
$env:SERVER_URL = "http://localhost:3000"

Write-Host "========================================="
Write-Host "Testing SQLBots Client"
Write-Host "========================================="
Write-Host "ENCRYPTION_KEY: Set"
Write-Host "SERVER_URL: $env:SERVER_URL"
Write-Host ""
Write-Host "Note: This will prompt for API Key"
Write-Host "You can enter a test API key to see the interface"
Write-Host "========================================="
Write-Host ""

# 运行客户端
cd $PSScriptRoot
go run main.go

