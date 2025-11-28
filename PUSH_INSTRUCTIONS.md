# 如何推送代码到 GitHub

## 方法 1：使用推送脚本（最简单）

1. 获取 GitHub Personal Access Token：
   - 访问：https://github.com/settings/tokens
   - 点击 "Generate new token (classic)"
   - 勾选 `repo` 权限
   - 点击 "Generate token"
   - **复制并保存 token**（只显示一次）

2. 运行推送脚本：
   ```bash
   cd /root/SQLBots
   ./push-to-github.sh YOUR_TOKEN_HERE
   ```

## 方法 2：使用 git 命令

```bash
cd /root/SQLBots

# 使用 token 设置远程 URL（临时）
git remote set-url origin https://YOUR_TOKEN@github.com/WongYIThong1/sqlbot-client-server.git

# 推送
git push -u origin main

# 恢复原始 URL（安全）
git remote set-url origin https://github.com/WongYIThong1/sqlbot-client-server.git
```

## 方法 3：配置 Git 凭据（永久）

```bash
cd /root/SQLBots

# 配置凭据存储
git config --global credential.helper store

# 推送（会提示输入用户名和密码）
# 用户名：你的 GitHub 用户名
# 密码：使用 Personal Access Token（不是 GitHub 密码）
git push -u origin main
```

## 当前状态

✅ 代码已提交到本地仓库  
✅ 远程仓库已配置  
⏳ 等待身份验证后推送

## 注意事项

- Personal Access Token 需要 `repo` 权限
- Token 只显示一次，请妥善保存
- 不要将 token 提交到代码仓库

