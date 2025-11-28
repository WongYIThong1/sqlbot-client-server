/**
 * 错误码定义
 */
export const ErrorCodes = {
  SUCCESS: 'SUCCESS',
  INVALID_API_KEY: 'INVALID_API_KEY',
  DECRYPTION_FAILED: 'DECRYPTION_FAILED',
  LICENSE_EXPIRED: 'LICENSE_EXPIRED',
  MACHINE_LIMIT_EXCEEDED: 'MACHINE_LIMIT_EXCEEDED',
  SERVER_ERROR: 'SERVER_ERROR',
};

/**
 * 创建错误响应
 * @param {string} statusCode - 错误码
 * @param {string} message - 错误消息（可选）
 * @returns {object} 错误响应对象
 */
export function createErrorResponse(statusCode, message = null) {
  const response = {
    status_code: statusCode,
  };
  
  if (message) {
    response.message = message;
  }
  
  return response;
}

/**
 * 创建成功响应
 * @param {object} data - 响应数据
 * @returns {object} 成功响应对象
 */
export function createSuccessResponse(data = {}) {
  return {
    status_code: ErrorCodes.SUCCESS,
    ...data,
  };
}

