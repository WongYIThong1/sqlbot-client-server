import {
  findMachineByMachineIdAndApiKey,
  countUserMachines,
  createMachine,
  updateMachine,
} from '../utils/database.js';
import { ErrorCodes } from '../utils/errors.js';

const MAX_MACHINES_PER_USER = 3;

/**
 * 验证或注册机器
 * @param {string} machineId - 机器 ID
 * @param {string} apiKey - API Key
 * @param {string} machineName - 机器名称
 * @param {number} ram - RAM 大小（GB）
 * @param {number} cores - CPU 核心数
 * @returns {Promise<object>} 机器信息和操作结果
 */
export async function verifyOrRegisterMachine(machineId, apiKey, machineName, ram, cores) {
  // 查找现有机器
  let machine = await findMachineByMachineIdAndApiKey(machineId, apiKey);
  
  if (!machine) {
    // 机器不存在，检查用户机器数量
    const machineCount = await countUserMachines(apiKey);
    
    if (machineCount >= MAX_MACHINES_PER_USER) {
      return {
        success: false,
        error: ErrorCodes.MACHINE_LIMIT_EXCEEDED,
        message: `Maximum ${MAX_MACHINES_PER_USER} machines allowed per user`,
      };
    }
    
    // 注册新机器
    machine = await createMachine({
      machine: machineId,
      api_key: apiKey,
      name: machineName,
      ram: ram,
      cores: cores,
    });
    
    return {
      success: true,
      machine: machine,
      isNew: true,
    };
  }
  
  // 机器已存在，返回机器信息
  return {
    success: true,
    machine: machine,
    isNew: false,
  };
}

/**
 * 更新机器信息
 * @param {string} machineId - 机器 ID
 * @param {string} apiKey - API Key
 * @param {object} updateData - 更新数据
 * @returns {Promise<object>} 更新后的机器对象
 */
export async function updateMachineInfo(machineId, apiKey, updateData) {
  return await updateMachine(machineId, apiKey, updateData);
}


