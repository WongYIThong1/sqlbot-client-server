import { decrypt, encrypt } from '../utils/encryption.js';
import { ErrorCodes, createErrorResponse, createSuccessResponse } from '../utils/errors.js';
import { verifyOrRegisterMachine } from '../services/machine.js';
import { verifyAndUpdateHardware } from '../services/hardware.js';
import { verifyLicense } from '../services/license.js';
import { getSessionKey, getOrCreateSessionKey } from '../utils/session.js';

/**
 * 心跳路由处理
 */
export async function heartbeatHandler(request, reply) {
  try {
    const { encrypted_data, use_session_key } = request.body;
    const user = request.user; // 从中间件获取
    const initialEncryptionKey = process.env.ENCRYPTION_KEY;
    
    // 确定使用哪个密钥：优先使用会话密钥，如果没有则使用初始密钥
    let encryptionKey = initialEncryptionKey;
    if (use_session_key !== false) {
      // 尝试获取会话密钥
      const sessionKey = getSessionKey(user.id);
      if (sessionKey) {
        encryptionKey = sessionKey;
      } else {
        // 会话密钥不存在或已过期，创建新的
        encryptionKey = getOrCreateSessionKey(user.id);
      }
    }
    
    if (!encrypted_data) {
      return reply.code(400).send(
        createErrorResponse(ErrorCodes.DECRYPTION_FAILED, 'encrypted_data is required')
      );
    }
    
    // 解密请求数据
    let decryptedData;
    try {
      const decryptedText = decrypt(encrypted_data, encryptionKey);
      decryptedData = JSON.parse(decryptedText);
    } catch (error) {
      // 如果使用会话密钥解密失败，尝试使用初始密钥（向后兼容）
      if (encryptionKey !== initialEncryptionKey) {
        try {
          const decryptedText = decrypt(encrypted_data, initialEncryptionKey);
          decryptedData = JSON.parse(decryptedText);
          // 解密成功，说明客户端还在使用初始密钥，清除会话密钥
          // 下次心跳时会创建新的会话密钥
        } catch (fallbackError) {
          return reply.code(400).send(
            createErrorResponse(ErrorCodes.DECRYPTION_FAILED, `Decryption failed: ${error.message}`)
          );
        }
      } else {
        return reply.code(400).send(
          createErrorResponse(ErrorCodes.DECRYPTION_FAILED, `Decryption failed: ${error.message}`)
        );
      }
    }
    
    const { machine_id, machine_name, ram, cores } = decryptedData;
    
    if (!machine_id || !machine_name || ram === undefined || cores === undefined) {
      return reply.code(400).send(
        createErrorResponse(ErrorCodes.SERVER_ERROR, 'Missing required fields in decrypted data')
      );
    }
    
    // 1. 验证或注册机器
    const machineResult = await verifyOrRegisterMachine(
      machine_id,
      user.api_key,
      machine_name,
      ram,
      cores
    );
    
    if (!machineResult.success) {
      return reply.code(403).send(
        createErrorResponse(machineResult.error, machineResult.message)
      );
    }
    
    // 2. 验证硬件信息（如果机器已存在）
    if (!machineResult.isNew) {
      await verifyAndUpdateHardware(
        machine_id,
        user.api_key,
        machineResult.machine,
        ram,
        cores
      );
    }
    
    // 3. 验证许可证
    const licenseResult = await verifyLicense(user.id);
    
    if (!licenseResult.valid) {
      return reply.code(403).send(
        createErrorResponse(licenseResult.error, licenseResult.message)
      );
    }
    
    // 构建响应数据
    const responseData = createSuccessResponse({
      license_info: licenseResult.licenseInfo,
      machine_info: {
        id: machineResult.machine.id,
        name: machineResult.machine.name,
        registered_at: machineResult.machine.created_at,
      },
    });
    
    // 加密响应数据
    const encryptedResponse = encrypt(JSON.stringify(responseData), encryptionKey);
    
    return reply.send({
      encrypted_data: encryptedResponse,
    });
    
  } catch (error) {
    console.error('Heartbeat error:', error);
    return reply.code(500).send(
      createErrorResponse(ErrorCodes.SERVER_ERROR, `Internal server error: ${error.message}`)
    );
  }
}

