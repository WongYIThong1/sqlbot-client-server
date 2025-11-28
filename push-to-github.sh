#!/bin/bash
# Script to push code to GitHub
# Usage: ./push-to-github.sh [GITHUB_TOKEN]

cd /root/SQLBots

if [ -z "$1" ]; then
    echo "Usage: ./push-to-github.sh <GITHUB_TOKEN>"
    echo ""
    echo "To get a GitHub Personal Access Token:"
    echo "1. Go to https://github.com/settings/tokens"
    echo "2. Click 'Generate new token (classic)'"
    echo "3. Select 'repo' scope"
    echo "4. Copy the token"
    echo ""
    echo "Then run: ./push-to-github.sh YOUR_TOKEN"
    exit 1
fi

GITHUB_TOKEN=$1
GITHUB_USER="WongYIThong1"
REPO_NAME="sqlbot-client-server"

# Update remote URL with token
git remote set-url origin https://${GITHUB_TOKEN}@github.com/${GITHUB_USER}/${REPO_NAME}.git

# Push to GitHub
echo "Pushing to GitHub..."
git push -u origin main

# Reset remote URL (remove token for security)
git remote set-url origin https://github.com/${GITHUB_USER}/${REPO_NAME}.git

echo "Done!"

