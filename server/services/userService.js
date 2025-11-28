const supabase = require('../db/supabase');

/**
 * Validate API key and get user info
 * @param {string} apiKey - API key to validate
 * @returns {Promise<object|null>} - User object or null if invalid
 */
async function validateApiKey(apiKey) {
    const { data, error } = await supabase
        .from('users')
        .select('*')
        .eq('api_key', apiKey)
        .single();

    if (error || !data) {
        return null;
    }
    return data;
}

module.exports = {
    validateApiKey,
};


