/**
 * Health check route
 * @param {object} fastify - Fastify instance
 */
async function healthRoutes(fastify) {
    fastify.get('/health', async (request, reply) => {
        return { 
            status: 'ok', 
            timestamp: new Date().toISOString() 
        };
    });
}

module.exports = healthRoutes;

