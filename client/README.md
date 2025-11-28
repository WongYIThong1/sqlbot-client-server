# SQLBots 客户端

基于 Go 的机器监控客户端，用于与 SQLBots 服务器通信。

## 功能特性

- 加密通信（AES-256-CBC）
- 自动机器注册
- 硬件信息收集
- 心跳机制（每 10 分钟）
- 优雅关闭

## 安装

```bash
go mod download
```

## 构建

```bash
go build -o sqlbots-client main.go
```

## 运行

### 使用环境变量（推荐）

设置环境变量后直接运行：

```bash
# Windows PowerShell
$env:API_KEY="your-api-key"
$env:ENCRYPTION_KEY="your-32-character-encryption-key"
./sqlbots-client

# Windows CMD
set API_KEY=your-api-key
set ENCRYPTION_KEY=your-32-character-encryption-key
sqlbots-client.exe

# Linux/Mac
export API_KEY="your-api-key"
export ENCRYPTION_KEY="your-32-character-encryption-key"
./sqlbots-client
```

### 使用命令行参数

```bash
./sqlbots-client \
  --api-key "your-api-key" \
  --server-url "https://api.sqlbots.online" \
  --encryption-key "your-32-character-encryption-key"
```

### 混合使用

命令行参数会覆盖环境变量。如果设置了环境变量，可以只传入必要的参数：

```bash
# 如果已设置 ENCRYPTION_KEY 环境变量，只需传入 API_KEY
./sqlbots-client --api-key "your-api-key"
```

## 配置参数

- `--api-key` / `API_KEY`: API Key（必填）
- `--server-url` / `SERVER_URL`: 服务器 URL（默认: https://api.sqlbots.online）
- `--encryption-key` / `ENCRYPTION_KEY`: 加密密钥（必填，32字符）

## 错误处理

客户端会根据服务器返回的错误码决定行为：

- **INVALID_API_KEY**: 终止程序
- **LICENSE_EXPIRED**: 终止程序
- **MACHINE_LIMIT_EXCEEDED**: 终止程序
- **网络错误**: 记录日志，继续运行

## 依赖

- `github.com/denisbrodbeck/machineid`: 获取机器唯一 ID
- `github.com/shirou/gopsutil/v3`: 获取系统硬件信息

