import { createClient } from '@supabase/supabase-js';
import dotenv from 'dotenv';

dotenv.config();

/**
 * 初始化 Supabase 客户端
 */
export const supabase = createClient(
  process.env.SUPABASE_URL,
  process.env.SUPABASE_KEY
);

/**
 * 通过 API Key 查找用户
 * @param {string} apiKey - API Key
 * @returns {Promise<object|null>} 用户对象或 null
 */
export async function findUserByApiKey(apiKey) {
  const { data, error } = await supabase
    .from('users')
    .select('*')
    .eq('api_key', apiKey)
    .single();
  
  if (error) {
    if (error.code === 'PGRST116') {
      // 未找到记录
      return null;
    }
    throw error;
  }
  
  return data;
}

/**
 * 通过 machine_id 和 api_key 查找机器
 * @param {string} machineId - 机器 ID
 * @param {string} apiKey - API Key
 * @returns {Promise<object|null>} 机器对象或 null
 */
export async function findMachineByMachineIdAndApiKey(machineId, apiKey) {
  const { data, error } = await supabase
    .from('machines')
    .select('*')
    .eq('machine', machineId)
    .eq('api_key', apiKey)
    .single();
  
  if (error) {
    if (error.code === 'PGRST116') {
      // 未找到记录
      return null;
    }
    throw error;
  }
  
  return data;
}

/**
 * 统计用户的机器数量
 * @param {string} apiKey - API Key
 * @returns {Promise<number>} 机器数量
 */
export async function countUserMachines(apiKey) {
  const { count, error } = await supabase
    .from('machines')
    .select('*', { count: 'exact', head: true })
    .eq('api_key', apiKey);
  
  if (error) {
    throw error;
  }
  
  return count || 0;
}

/**
 * 创建新机器
 * @param {object} machineData - 机器数据
 * @returns {Promise<object>} 创建的机器对象
 */
export async function createMachine(machineData) {
  const { data, error } = await supabase
    .from('machines')
    .insert(machineData)
    .select()
    .single();
  
  if (error) {
    throw error;
  }
  
  return data;
}

/**
 * 更新机器信息
 * @param {string} machineId - 机器 ID
 * @param {string} apiKey - API Key
 * @param {object} updateData - 更新数据
 * @returns {Promise<object>} 更新后的机器对象
 */
export async function updateMachine(machineId, apiKey, updateData) {
  const { data, error } = await supabase
    .from('machines')
    .update({
      ...updateData,
      updated_at: new Date().toISOString(),
    })
    .eq('machine', machineId)
    .eq('api_key', apiKey)
    .select()
    .single();
  
  if (error) {
    throw error;
  }
  
  return data;
}

/**
 * 通过用户 ID 查找许可证
 * @param {string} userId - 用户 ID
 * @returns {Promise<object|null>} 许可证对象或 null
 */
export async function findLicenseByUserId(userId) {
  const { data, error } = await supabase
    .from('licenses')
    .select('*')
    .eq('user_id', userId)
    .single();
  
  if (error) {
    if (error.code === 'PGRST116') {
      // 未找到记录
      return null;
    }
    throw error;
  }
  
  return data;
}


