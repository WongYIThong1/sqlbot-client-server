const CryptoJS = require('crypto-js');
const config = require('../config');

const ENCRYPTION_KEY = config.encryption.key;

/**
 * Decrypt heartbeat data
 * CryptoJS uses OpenSSL-compatible format: Salt (8 bytes) + IV (16 bytes) + Ciphertext
 * @param {string} encryptedData - Base64 encoded encrypted data
 * @returns {object|null} - Decrypted data object or null if decryption fails
 */
function decryptData(encryptedData) {
    try {
        // CryptoJS automatically handles salt and IV extraction
        const bytes = CryptoJS.AES.decrypt(encryptedData, ENCRYPTION_KEY);
        const decrypted = bytes.toString(CryptoJS.enc.Utf8);
        if (!decrypted) {
            return null;
        }
        return JSON.parse(decrypted);
    } catch (error) {
        console.error('Decryption error:', error);
        return null;
    }
}

/**
 * Encrypt response data
 * CryptoJS automatically generates salt and IV
 * @param {object} data - Data object to encrypt
 * @returns {string|null} - Base64 encoded encrypted string or null if encryption fails
 */
function encryptData(data) {
    try {
        const encrypted = CryptoJS.AES.encrypt(JSON.stringify(data), ENCRYPTION_KEY);
        return encrypted.toString();
    } catch (error) {
        console.error('Encryption error:', error);
        return null;
    }
}

module.exports = {
    decryptData,
    encryptData,
};


