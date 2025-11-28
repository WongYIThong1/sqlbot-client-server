const supabase = require('../db/supabase');

/**
 * Get machine count for API key
 * @param {string} apiKey - API key
 * @returns {Promise<number>} - Number of machines registered
 */
async function getMachineCount(apiKey) {
    const { data, error } = await supabase
        .from('machines')
        .select('*')
        .eq('api_key', apiKey);

    if (error) {
        return 0;
    }
    return data ? data.length : 0;
}

/**
 * Get machine by API key and machine identifier
 * @param {string} apiKey - API key
 * @param {string} machineId - Machine identifier
 * @returns {Promise<object|null>} - Machine object or null if not found
 */
async function getMachineByIdentifier(apiKey, machineId) {
    const { data, error } = await supabase
        .from('machines')
        .select('*')
        .eq('api_key', apiKey)
        .eq('machine', machineId)
        .single();

    if (error || !data) {
        return null;
    }
    return data;
}

/**
 * Register new machine
 * @param {string} apiKey - API key
 * @param {string} name - Machine name
 * @param {string} machineId - Machine identifier
 * @param {number} ram - RAM in GB
 * @param {number} cores - CPU cores
 * @returns {Promise<object|null>} - Registered machine object or null if failed
 */
async function registerMachine(apiKey, name, machineId, ram, cores) {
    const { data, error } = await supabase
        .from('machines')
        .insert([
            {
                name: name,
                machine: machineId,
                api_key: apiKey,
                ram: ram,
                cores: cores
            }
        ])
        .select()
        .single();

    if (error) {
        console.error('Error registering machine:', error);
        return null;
    }
    return data;
}

/**
 * Update machine info
 * @param {string} apiKey - API key
 * @param {string} machineId - Machine identifier
 * @param {number} ram - RAM in GB
 * @param {number} cores - CPU cores
 * @returns {Promise<object|null>} - Updated machine object or null if failed
 */
async function updateMachine(apiKey, machineId, ram, cores) {
    const { data, error } = await supabase
        .from('machines')
        .update({ ram: ram, cores: cores })
        .eq('api_key', apiKey)
        .eq('machine', machineId)
        .select()
        .single();

    if (error) {
        console.error('Error updating machine:', error);
        return null;
    }
    return data;
}

/**
 * Verify machine hardware matches stored data
 * @param {object} machine - Machine object from database
 * @param {number} ram - Current RAM in GB
 * @param {number} cores - Current CPU cores
 * @returns {boolean} - True if hardware matches
 */
function verifyHardware(machine, ram, cores) {
    return machine.ram === ram && machine.cores === cores;
}

module.exports = {
    getMachineCount,
    getMachineByIdentifier,
    registerMachine,
    updateMachine,
    verifyHardware,
};


