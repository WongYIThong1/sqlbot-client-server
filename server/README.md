# SQLBots 服务器端

基于 Node.js + Fastify + Supabase 的机器监控与管理系统服务器端。

## 功能特性

- 加密通信（AES-256-CBC）
- API Key 认证
- 许可证验证
- 机器注册与管理
- 硬件信息验证
- 心跳机制

## 安装

```bash
npm install
```

## 配置

复制 `.env.example` 为 `.env` 并填写配置：

```bash
cp .env.example .env
```

## 运行

```bash
npm start
```

开发模式（自动重启）：

```bash
npm run dev
```

## 环境变量

- `SUPABASE_URL`: Supabase 项目 URL
- `SUPABASE_KEY`: Supabase 匿名密钥
- `ENCRYPTION_KEY`: 加密密钥（32字符）
- `PORT`: 服务器端口（默认: 3000）
- `HOST`: 服务器主机（默认: 0.0.0.0）


