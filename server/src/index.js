import Fastify from 'fastify';
import dotenv from 'dotenv';
import { apiKeyAuth } from './middleware/auth.js';
import { heartbeatHandler } from './routes/heartbeat.js';
import { keyExchangeHandler } from './routes/keyExchange.js';

dotenv.config();

const fastify = Fastify({
  logger: true,
});

// 密钥交换路由
fastify.post('/key-exchange', async (request, reply) => {
  // 先执行 API Key 验证
  const authResult = await apiKeyAuth(request, reply);
  if (authResult) {
    return authResult;
  }
  // 验证通过，处理密钥交换
  return keyExchangeHandler(request, reply);
});

// 注册路由（中间件在路由处理函数中调用）
fastify.post('/heartbeat', async (request, reply) => {
  // 先执行 API Key 验证
  const authResult = await apiKeyAuth(request, reply);
  if (authResult) {
    // 如果返回了响应，说明验证失败，直接返回
    return authResult;
  }
  // 验证通过，继续处理心跳
  return heartbeatHandler(request, reply);
});

// 健康检查路由
fastify.get('/health', async (request, reply) => {
  return { status: 'ok' };
});

// 启动服务器
const start = async () => {
  try {
    const host = process.env.HOST || '0.0.0.0';
    const port = parseInt(process.env.PORT || '3000', 10);
    
    await fastify.listen({ host, port });
    console.log(`Server listening on ${host}:${port}`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();

