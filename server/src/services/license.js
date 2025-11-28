import { findLicenseByUserId } from '../utils/database.js';
import { ErrorCodes } from '../utils/errors.js';

/**
 * 验证许可证
 * @param {string} userId - 用户 ID
 * @returns {Promise<object>} 验证结果和许可证信息
 */
export async function verifyLicense(userId) {
  const license = await findLicenseByUserId(userId);
  
  if (!license) {
    return {
      valid: false,
      error: ErrorCodes.LICENSE_EXPIRED,
      message: 'No license found for user',
    };
  }
  
  // 检查许可证是否过期
  if (license.expires_at) {
    const expiresAt = new Date(license.expires_at);
    const now = new Date();
    
    if (now > expiresAt) {
      return {
        valid: false,
        error: ErrorCodes.LICENSE_EXPIRED,
        message: 'License has expired',
        expiresAt: license.expires_at,
      };
    }
  }
  
  return {
    valid: true,
    licenseInfo: {
      expires_at: license.expires_at,
      plan_type: license.plan_type,
    },
  };
}


