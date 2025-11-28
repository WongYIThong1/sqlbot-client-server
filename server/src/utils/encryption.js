import CryptoJS from 'crypto-js';

/**
 * EVP_BytesToKey 密钥派生（兼容 OpenSSL 和 CryptoJS）
 * @param {string} password - 密码
 * @param {string} salt - 盐值
 * @returns {object} 包含 key 和 iv 的对象
 */
function evpBytesToKey(password, salt) {
  const keySize = 256 / 32; // 32 bytes for AES-256
  const ivSize = 128 / 32;  // 16 bytes for IV
  const iterations = 1;
  
  let derivedKey = CryptoJS.enc.Hex.parse('');
  let derivedIV = CryptoJS.enc.Hex.parse('');
  let derived = CryptoJS.enc.Hex.parse('');
  
  // 第一次迭代：MD5(password + salt)
  derived = CryptoJS.MD5(CryptoJS.enc.Utf8.parse(password).concat(salt));
  
  // 生成 key
  derivedKey = derived;
  for (let i = 1; i < keySize; i++) {
    derived = CryptoJS.MD5(derived.concat(CryptoJS.enc.Utf8.parse(password)).concat(salt));
    derivedKey = derivedKey.concat(derived);
  }
  
  // 生成 IV
  derived = CryptoJS.MD5(derivedKey.concat(CryptoJS.enc.Utf8.parse(password)).concat(salt));
  derivedIV = derived;
  for (let i = 1; i < ivSize; i++) {
    derived = CryptoJS.MD5(derived.concat(CryptoJS.enc.Utf8.parse(password)).concat(salt));
    derivedIV = derivedIV.concat(derived);
  }
  
  return {
    key: CryptoJS.lib.WordArray.create(derivedKey.words.slice(0, keySize * 4)),
    iv: CryptoJS.lib.WordArray.create(derivedIV.words.slice(0, ivSize * 4)),
  };
}

/**
 * 加密数据（OpenSSL 兼容格式）
 * @param {string} plaintext - 明文
 * @param {string} password - 密码（加密密钥）
 * @returns {string} Base64 编码的加密数据
 */
export function encrypt(plaintext, password) {
  // 生成随机盐值（8字节）
  const salt = CryptoJS.lib.WordArray.random(64 / 8);
  
  // 派生密钥和 IV
  const { key, iv } = evpBytesToKey(password, salt);
  
  // 加密
  const encrypted = CryptoJS.AES.encrypt(plaintext, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });
  
  // OpenSSL 格式：Salted__ + salt (8 bytes) + encrypted data
  const saltedPrefix = CryptoJS.enc.Utf8.parse('Salted__');
  const combined = saltedPrefix.concat(salt).concat(encrypted.ciphertext);
  
  // Base64 编码
  return CryptoJS.enc.Base64.stringify(combined);
}

/**
 * 解密数据（OpenSSL 兼容格式）
 * @param {string} ciphertext - Base64 编码的加密数据
 * @param {string} password - 密码（加密密钥）
 * @returns {string} 解密后的明文
 */
export function decrypt(ciphertext, password) {
  try {
    // Base64 解码
    const encryptedData = CryptoJS.enc.Base64.parse(ciphertext);
    
    // 检查格式：前8字节应该是 "Salted__"
    const saltedPrefix = CryptoJS.enc.Utf8.parse('Salted__');
    const prefixBytes = CryptoJS.lib.WordArray.create(encryptedData.words.slice(0, 2));
    
    if (saltedPrefix.toString() !== prefixBytes.toString()) {
      throw new Error('Invalid encrypted data format');
    }
    
    // 提取 salt（接下来的8字节）
    const salt = CryptoJS.lib.WordArray.create(encryptedData.words.slice(2, 4));
    
    // 提取加密数据（剩余部分）
    const ciphertextOnly = CryptoJS.lib.WordArray.create(encryptedData.words.slice(4));
    
    // 派生密钥和 IV
    const { key, iv } = evpBytesToKey(password, salt);
    
    // 创建加密参数对象
    const cipherParams = CryptoJS.lib.CipherParams.create({
      ciphertext: ciphertextOnly,
    });
    
    // 解密
    const decrypted = CryptoJS.AES.decrypt(cipherParams, key, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    
    return decrypted.toString(CryptoJS.enc.Utf8);
  } catch (error) {
    throw new Error(`Decryption failed: ${error.message}`);
  }
}

