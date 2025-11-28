/**
 * 会话密钥管理器
 * 管理每个用户的会话密钥
 */

// 内存存储会话密钥（生产环境应使用 Redis 等）
const sessionKeys = new Map();

// 会话密钥有效期（30分钟）
const SESSION_KEY_TTL = 30 * 60 * 1000;

/**
 * 生成随机会话密钥（32字符）
 */
function generateSessionKey() {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let key = '';
  for (let i = 0; i < 32; i++) {
    key += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return key;
}

/**
 * 获取或创建用户的会话密钥
 * @param {string} userId - 用户 ID
 * @returns {string} 会话密钥
 */
export function getOrCreateSessionKey(userId) {
  const session = sessionKeys.get(userId);
  
  // 如果会话存在且未过期，返回现有密钥
  if (session && Date.now() < session.expiresAt) {
    return session.key;
  }
  
  // 生成新会话密钥
  const newKey = generateSessionKey();
  sessionKeys.set(userId, {
    key: newKey,
    expiresAt: Date.now() + SESSION_KEY_TTL,
    createdAt: Date.now(),
  });
  
  return newKey;
}

/**
 * 获取用户的会话密钥（不创建新密钥）
 * @param {string} userId - 用户 ID
 * @returns {string|null} 会话密钥或 null
 */
export function getSessionKey(userId) {
  const session = sessionKeys.get(userId);
  
  if (session && Date.now() < session.expiresAt) {
    return session.key;
  }
  
  return null;
}

/**
 * 删除用户的会话密钥
 * @param {string} userId - 用户 ID
 */
export function deleteSessionKey(userId) {
  sessionKeys.delete(userId);
}

/**
 * 清理过期的会话密钥
 */
export function cleanupExpiredSessions() {
  const now = Date.now();
  for (const [userId, session] of sessionKeys.entries()) {
    if (now >= session.expiresAt) {
      sessionKeys.delete(userId);
    }
  }
}

// 定期清理过期会话（每5分钟）
setInterval(cleanupExpiredSessions, 5 * 60 * 1000);


