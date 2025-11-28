# SQLBots 项目

基于客户端-服务器架构的机器监控与管理系统。

## 项目结构

```
sqlbot-client-server/
├── server/          # 服务器端 (Node.js + Fastify + Supabase)
├── client/          # 客户端 (Go)
└── README.md        # 本文件
```

## 功能特性

### 核心功能

1. **加密通信**
   - AES-256-CBC 加密
   - OpenSSL 兼容格式
   - 双向加密（请求和响应）

2. **许可证管理**
   - 每次心跳验证许可证有效性
   - 过期自动检测

3. **机器注册与管理**
   - 首次心跳自动注册
   - 每用户最多 3 台机器
   - 存储机器名称、RAM、CPU 核心数

4. **硬件验证**
   - 每次心跳收集并验证硬件信息
   - 自动检测硬件配置变更

5. **心跳机制**
   - 每 10 分钟发送心跳
   - 启动时立即发送首次心跳

## 快速开始

### 服务器端

1. 进入服务器目录：
```bash
cd server
```

2. 安装依赖：
```bash
npm install
```

3. 配置环境变量：
```bash
cp .env.example .env
# 编辑 .env 文件，填写 Supabase 配置和加密密钥
```

4. 启动服务器：
```bash
npm start
```

### 客户端

1. 进入客户端目录：
```bash
cd client
```

2. 安装依赖：
```bash
go mod download
```

3. 构建：
```bash
go build -o sqlbots-client main.go
```

4. 运行：
```bash
./sqlbots-client \
  --api-key "your-api-key" \
  --encryption-key "your-32-character-encryption-key"
```

## 配置说明

### 服务器端环境变量

- `SUPABASE_URL`: Supabase 项目 URL
- `SUPABASE_KEY`: Supabase 匿名密钥
- `ENCRYPTION_KEY`: 加密密钥（32字符）
- `PORT`: 服务器端口（默认: 3000）
- `HOST`: 服务器主机（默认: 0.0.0.0）

### 客户端配置

- `API_KEY`: API Key（必填）
- `SERVER_URL`: 服务器 URL（默认: https://api.sqlbots.online）
- `ENCRYPTION_KEY`: 加密密钥（必填，必须与服务器一致）

## 安全特性

1. API Key 认证
2. AES-256-CBC 端到端加密
3. 硬件信息验证防篡改
4. 严格许可证控制
5. 每用户最多 3 台机器限制

## 技术栈

### 服务器端
- Node.js
- Fastify
- Supabase (PostgreSQL)
- crypto-js

### 客户端
- Go 1.21+
- gopsutil (系统信息)
- machineid (机器唯一标识)

## 许可证

ISC

