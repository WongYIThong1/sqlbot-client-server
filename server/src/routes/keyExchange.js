import { encrypt } from '../utils/encryption.js';
import { ErrorCodes, createErrorResponse, createSuccessResponse } from '../utils/errors.js';
import { getOrCreateSessionKey } from '../utils/session.js';

/**
 * 密钥交换路由处理
 * 客户端使用初始 ENCRYPTION_KEY 请求会话密钥
 */
export async function keyExchangeHandler(request, reply) {
  try {
    const user = request.user; // 从中间件获取（已通过 API Key 验证）
    const initialEncryptionKey = process.env.ENCRYPTION_KEY;
    
    // 获取或创建会话密钥
    const sessionKey = getOrCreateSessionKey(user.id);
    
    // 使用初始密钥加密会话密钥
    const encryptedSessionKey = encrypt(sessionKey, initialEncryptionKey);
    
    // 构建响应（包含用户名信息）
    const responseData = createSuccessResponse({
      session_key: encryptedSessionKey,
      expires_in: 1800, // 30分钟（秒）
      username: user.username || user.email || 'User', // 返回用户名
    });
    
    return reply.send(responseData);
    
  } catch (error) {
    console.error('Key exchange error:', error);
    return reply.code(500).send(
      createErrorResponse(ErrorCodes.SERVER_ERROR, `Internal server error: ${error.message}`)
    );
  }
}

