const encryptionService = require('../services/encryption');
const userService = require('../services/userService');
const machineService = require('../services/machineService');
const licenseService = require('../services/licenseService');

const MAX_MACHINES = 3;

/**
 * Heartbeat route handler
 * @param {object} fastify - Fastify instance
 */
async function heartbeatRoutes(fastify) {
    fastify.post('/heartbeat', async (request, reply) => {
        try {
            const { data: encryptedData } = request.body;

            // Validate encrypted data presence
            if (!encryptedData) {
                return reply.status(400).send({
                    success: false,
                    error: 'Missing encrypted data'
                });
            }

            // Decrypt the heartbeat data
            const heartbeatData = encryptionService.decryptData(encryptedData);

            if (!heartbeatData) {
                return reply.status(400).send({
                    success: false,
                    error: 'Invalid encrypted data'
                });
            }

            // Validate required fields
            if (!heartbeatData.apiKey || !heartbeatData.machineName || !heartbeatData.machineId || 
                heartbeatData.ram === undefined || heartbeatData.cores === undefined) {
                return reply.status(400).send({
                    success: false,
                    error: 'Missing required fields in heartbeat data'
                });
            }

            const { apiKey, machineName, machineId, ram, cores } = heartbeatData;

            // Validate API key
            const user = await userService.validateApiKey(apiKey);

            if (!user) {
                const responseData = encryptionService.encryptData({
                    success: false,
                    error: 'Invalid API key',
                    code: 'INVALID_API_KEY'
                });
                return reply.status(401).send({ data: responseData });
            }

            // Check license expiration
            const licenseValidation = await licenseService.validateLicense(user.license_id);

            if (!licenseValidation.valid) {
                const responseData = encryptionService.encryptData({
                    success: false,
                    error: 'License has expired',
                    code: 'LICENSE_EXPIRED'
                });
                return reply.status(403).send({ data: responseData });
            }

            // Check if machine exists
            const existingMachine = await machineService.getMachineByIdentifier(apiKey, machineId);

            if (!existingMachine) {
                // Check machine limit (max 3 machines)
                const machineCount = await machineService.getMachineCount(apiKey);

                if (machineCount >= MAX_MACHINES) {
                    const responseData = encryptionService.encryptData({
                        success: false,
                        error: 'Maximum machine limit reached. Unable to configure new machine.',
                        code: 'MACHINE_LIMIT_EXCEEDED'
                    });
                    return reply.status(403).send({ data: responseData });
                }

                // Register new machine
                const newMachine = await machineService.registerMachine(apiKey, machineName, machineId, ram, cores);

                if (!newMachine) {
                    return reply.status(500).send({
                        success: false,
                        error: 'Failed to register machine'
                    });
                }

                const responseData = encryptionService.encryptData({
                    success: true,
                    message: 'Machine registered successfully',
                    isNewMachine: true,
                    licenseValidUntil: licenseValidation.expiresAt
                });

                return reply.send({ data: responseData });
            }

            // Machine exists - verify RAM and Cores
            const hardwareMatches = machineService.verifyHardware(existingMachine, ram, cores);

            if (!hardwareMatches) {
                // Hardware mismatch - update machine info
                await machineService.updateMachine(apiKey, machineId, ram, cores);

                const responseData = encryptionService.encryptData({
                    success: true,
                    message: 'Machine hardware info updated',
                    hardwareChanged: true,
                    licenseValidUntil: licenseValidation.expiresAt
                });

                return reply.send({ data: responseData });
            }

            // Everything matches - successful heartbeat
            const responseData = encryptionService.encryptData({
                success: true,
                message: 'Heartbeat received',
                licenseValidUntil: licenseValidation.expiresAt
            });

            return reply.send({ data: responseData });

        } catch (error) {
            console.error('Heartbeat error:', error);
            return reply.status(500).send({
                success: false,
                error: 'Internal server error'
            });
        }
    });
}

module.exports = heartbeatRoutes;


