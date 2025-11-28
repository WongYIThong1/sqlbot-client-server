const fastify = require('fastify')({ logger: true });
const config = require('./config');

// Register routes
const healthRoutes = require('./routes/health');
const heartbeatRoutes = require('./routes/heartbeat');

// Register route plugins
fastify.register(healthRoutes);
fastify.register(heartbeatRoutes);

// Start server
const start = async () => {
    try {
        const { port, host } = config.server;

        await fastify.listen({ port, host });
        console.log(`Server running on ${host}:${port}`);
    } catch (err) {
        fastify.log.error(err);
        process.exit(1);
    }
};

start();
