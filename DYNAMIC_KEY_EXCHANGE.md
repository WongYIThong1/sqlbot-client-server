# 动态密钥交换功能说明

## 功能概述

实现了动态密钥交换机制，每次客户端连接时都会生成新的会话密钥，提高系统安全性。

## 工作流程

### 1. 密钥交换流程

```
客户端                          服务器
  │                              │
  │  1. 使用 API_KEY 请求会话密钥 │
  │  POST /key-exchange          │
  │  { API_KEY: "xxx" } ────────→│
  │                              │
  │                              │  2. 验证 API_KEY
  │                              │  3. 生成会话密钥（32字符随机字符串）
  │                              │  4. 使用初始 ENCRYPTION_KEY 加密会话密钥
  │                              │
  │  ←───────────────────────────│
  │  {                            │
  │    status_code: "SUCCESS",    │
  │    session_key: "加密的密钥",  │
  │    expires_in: 1800          │
  │  }                            │
  │                              │
  │  5. 使用初始密钥解密会话密钥   │
  │  6. 保存会话密钥（30分钟有效） │
```

### 2. 心跳通信流程

```
客户端                          服务器
  │                              │
  │  1. 检查是否有有效会话密钥     │
  │  2. 使用会话密钥加密数据       │
  │  POST /heartbeat             │
  │  {                            │
  │    API_KEY: "xxx",           │
  │    encrypted_data: "...",    │
  │    use_session_key: true     │
  │  } ──────────────────────────→│
  │                              │
  │                              │  3. 验证 API_KEY
  │                              │  4. 获取用户的会话密钥
  │                              │  5. 使用会话密钥解密
  │                              │  6. 处理业务逻辑
  │                              │  7. 使用会话密钥加密响应
  │                              │
  │  ←───────────────────────────│
  │  { encrypted_data: "..." }  │
  │                              │
  │  8. 使用会话密钥解密响应       │
```

### 3. 密钥刷新机制

- **自动刷新**：会话密钥过期前，客户端会在下次心跳时自动交换新密钥
- **过期时间**：30分钟（1800秒）
- **向后兼容**：如果会话密钥解密失败，服务器会尝试使用初始密钥

## 安全优势

### 1. 动态密钥
- 每次连接使用不同的会话密钥
- 即使某个会话密钥泄露，也只影响该会话
- 定期自动刷新，减少密钥暴露时间

### 2. 多层保护
- **初始密钥**：用于密钥交换（长期有效）
- **会话密钥**：用于实际通信（短期有效，30分钟）
- **API_KEY**：身份验证（用户级别）

### 3. 向后兼容
- 如果客户端没有会话密钥，可以使用初始密钥
- 服务器支持两种密钥的解密
- 平滑升级，不影响现有客户端

## 代码结构

### 服务器端

- `server/src/utils/session.js` - 会话密钥管理器
- `server/src/routes/keyExchange.js` - 密钥交换端点
- `server/src/routes/heartbeat.js` - 更新支持会话密钥

### 客户端

- `client/session/manager.go` - 会话密钥管理器
- `client/keyexchange/exchange.go` - 密钥交换逻辑
- `client/heartbeat/sender.go` - 更新支持会话密钥
- `client/main.go` - 集成密钥交换流程

## 配置说明

### 服务器端

无需额外配置，使用现有的 `ENCRYPTION_KEY` 作为初始密钥。

### 客户端

无需额外配置，自动从环境变量或命令行参数读取 `ENCRYPTION_KEY`。

## 使用示例

### 客户端启动流程

1. **启动时**：
   ```
   Exchanging session key...
   Session key exchanged successfully
   Sending initial heartbeat...
   Initial heartbeat successful
   ```

2. **定期心跳**：
   ```
   Sending periodic heartbeat...
   Heartbeat successful
   ```

3. **密钥刷新**：
   ```
   Session key expired or missing, exchanging new key...
   Session key refreshed
   Sending periodic heartbeat...
   ```

## 注意事项

1. **会话密钥存储**：服务器端使用内存存储（生产环境建议使用 Redis）
2. **密钥过期**：会话密钥30分钟后自动过期，客户端会自动刷新
3. **错误处理**：如果密钥交换失败，客户端会回退使用初始密钥
4. **性能影响**：密钥交换只在启动和过期时进行，不影响正常心跳性能

## 未来改进

1. **Redis 存储**：使用 Redis 存储会话密钥，支持多服务器部署
2. **密钥轮换**：定期自动轮换初始密钥
3. **密钥预取**：在密钥即将过期前提前获取新密钥
4. **监控告警**：监控密钥交换失败率和异常情况

