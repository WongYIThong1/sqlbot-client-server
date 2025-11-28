import { updateMachineInfo } from './machine.js';

/**
 * 验证硬件信息并更新（如果需要）
 * @param {string} machineId - 机器 ID
 * @param {string} apiKey - API Key
 * @param {object} currentMachine - 当前机器记录
 * @param {number} reportedRam - 客户端报告的 RAM（GB）
 * @param {number} reportedCores - 客户端报告的 CPU 核心数
 * @returns {Promise<object>} 验证结果和更新信息
 */
export async function verifyAndUpdateHardware(
  machineId,
  apiKey,
  currentMachine,
  reportedRam,
  reportedCores
) {
  const changes = {};
  const hasChanges = 
    currentMachine.ram !== reportedRam ||
    currentMachine.cores !== reportedCores;
  
  if (hasChanges) {
    if (currentMachine.ram !== reportedRam) {
      changes.ram = reportedRam;
    }
    if (currentMachine.cores !== reportedCores) {
      changes.cores = reportedCores;
    }
    
    // 更新数据库
    const updatedMachine = await updateMachineInfo(machineId, apiKey, changes);
    
    return {
      hasChanges: true,
      changes: changes,
      machine: updatedMachine,
    };
  }
  
  return {
    hasChanges: false,
    machine: currentMachine,
  };
}

