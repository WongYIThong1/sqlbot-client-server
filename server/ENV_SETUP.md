# 环境变量配置指南

## 如何配置 .env 文件

在 `server` 目录下创建一个名为 `.env` 的文件，并按照以下说明配置：

### 必需的环境变量

```env
# Supabase 配置
SUPABASE_URL=your_supabase_url_here
SUPABASE_KEY=your_supabase_anon_key_here

# 加密密钥（必须与客户端使用的密钥匹配）
# 注意：密钥长度建议为32个字符（用于AES-256）
ENCRYPTION_KEY=your-secret-encryption-key-32ch

# 服务器配置（可选，有默认值）
PORT=3000
HOST=0.0.0.0
```

## 详细配置说明

### 1. SUPABASE_URL (第9行)
这是你的 Supabase 项目 URL。

**如何获取：**
1. 登录到 [Supabase Dashboard](https://app.supabase.com/)
2. 选择你的项目
3. 进入 **Settings** → **API**
4. 在 **Project URL** 部分，复制完整的 URL
   - 格式类似：`https://xxxxxxxxxxxxx.supabase.co`

**示例：**
```env
SUPABASE_URL=https://abcdefghijklmnop.supabase.co
```

### 2. SUPABASE_KEY (第10行)
这是你的 Supabase 匿名/公开密钥（anon key）。

**如何获取：**
1. 在 Supabase Dashboard 中，进入 **Settings** → **API**
2. 在 **Project API keys** 部分
3. 找到 **anon** `public` 密钥
4. 点击复制按钮复制密钥
   - 这是一个很长的字符串，以 `eyJ` 开头

**示例：**
```env
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFiY2RlZmdoaWprbG1ub3AiLCJyb2xlIjoiYW5vbiIsImlhdCI6MTYxNjIzOTAyMiwiZXhwIjoxOTMxODE1MDIyfQ.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## 完整配置示例

创建一个 `.env` 文件（注意：不要提交到 Git，已在 .gitignore 中）：

```env
# Supabase Configuration
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InlvdXItcHJvamVjdC1pZCIsInJvbGUiOiJhbm9uIiwiaWF0IjoxNjE2MjM5MDIyLCJleHAiOjE5MzE4MTUwMjJ9.your-anon-key-here

# Encryption Key (must match client ENCRYPTION_KEY)
# Generate a secure 32-character key for production
ENCRYPTION_KEY=your-secret-encryption-key-32ch

# Server Configuration (optional)
PORT=3000
HOST=0.0.0.0
```

## 安全注意事项

1. **不要将 `.env` 文件提交到 Git**
   - 确保 `.env` 在 `.gitignore` 中
   - 只提交 `.env.example` 作为模板

2. **生产环境密钥**
   - 使用强随机生成的 `ENCRYPTION_KEY`
   - 可以使用以下命令生成：
     ```bash
     node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
     ```

3. **Supabase 密钥**
   - `SUPABASE_KEY` 是公开的 anon key，可以安全地在前端使用
   - 但确保已正确配置 Row Level Security (RLS) 策略

## 验证配置

配置完成后，启动服务器：
```bash
npm start
```

如果配置正确，服务器会启动并显示：
```
Server running on 0.0.0.0:3000
```

如果出现错误，检查：
- `.env` 文件是否在 `server` 目录下
- 所有必需的环境变量是否都已设置
- Supabase URL 和 Key 是否正确

