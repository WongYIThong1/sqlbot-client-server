const supabase = require('../db/supabase');

/**
 * Check if license is expired
 * @param {string} expirationDate - License expiration date string
 * @returns {boolean} - True if expired, false if valid
 */
function isLicenseExpired(expirationDate) {
    if (!expirationDate) return true;
    const expDate = new Date(expirationDate);
    const now = new Date();
    return now > expDate;
}

/**
 * Get license expiration date by license ID
 * @param {string} licenseId - License UUID
 * @returns {Promise<string|null>} - Expiration date or null
 */
async function getLicenseExpiration(licenseId) {
    if (!licenseId) {
        return null;
    }

    const { data: license, error } = await supabase
        .from('licenses')
        .select('expires_at')
        .eq('id', licenseId)
        .single();

    if (error || !license) {
        return null;
    }

    return license.expires_at;
}

/**
 * Check if user's license is valid
 * @param {string} licenseId - License UUID
 * @returns {Promise<{valid: boolean, expiresAt: string|null}>} - License validation result
 */
async function validateLicense(licenseId) {
    const expiresAt = await getLicenseExpiration(licenseId);
    
    if (!expiresAt) {
        return { valid: false, expiresAt: null };
    }

    const expired = isLicenseExpired(expiresAt);
    return { valid: !expired, expiresAt };
}

module.exports = {
    isLicenseExpired,
    getLicenseExpiration,
    validateLicense,
};


