import { findUserByApiKey } from '../utils/database.js';
import { ErrorCodes, createErrorResponse } from '../utils/errors.js';

/**
 * API Key 验证中间件
 */
export async function apiKeyAuth(request, reply) {
  const apiKey = request.body?.API_KEY;
  
  if (!apiKey) {
    return reply.code(400).send(
      createErrorResponse(ErrorCodes.INVALID_API_KEY, 'API_KEY is required')
    );
  }
  
  // 验证 API Key
  const user = await findUserByApiKey(apiKey);
  
  if (!user) {
    return reply.code(401).send(
      createErrorResponse(ErrorCodes.INVALID_API_KEY, 'Invalid API Key')
    );
  }
  
  // 将用户信息附加到请求对象，供后续使用
  request.user = user;
  
  return;
}


