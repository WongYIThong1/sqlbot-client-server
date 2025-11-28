require('dotenv').config();

module.exports = {
    supabase: {
        url: process.env.SUPABASE_URL,
        key: process.env.SUPABASE_KEY,
    },
    encryption: {
        key: process.env.ENCRYPTION_KEY || 'your-secret-encryption-key-32ch',
    },
    server: {
        port: process.env.PORT || 3000,
        host: process.env.HOST || '0.0.0.0',
    },
};


