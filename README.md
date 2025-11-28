# SQLBots - Machine Monitoring System

A client-server system for monitoring and managing machines with license validation and heartbeat functionality.

## Project Structure

```
SQLBots/
├── server/          # Node.js backend server (Fastify)
│   ├── index.js     # Main server file
│   ├── package.json # Dependencies
│   └── .env.example # Environment variables template
└── client/          # Go client (Windows)
    ├── main.go      # Client application
    └── go.mod       # Go dependencies
```

## Features

- **API Key Authentication**: All requests require valid API keys from Supabase
- **Machine Registration**: Automatic machine registration on first heartbeat
- **Hardware Verification**: Validates RAM and CPU cores on each heartbeat
- **License Management**: Checks license expiration and enforces limits
- **Machine Limit**: Maximum 3 machines per user
- **Encrypted Communication**: All heartbeat data is encrypted using AES
- **Periodic Heartbeats**: Sends heartbeat every 10 minutes

## Server Setup

### Prerequisites
- Node.js (v14 or higher)
- Supabase account with configured database

### Installation

1. Navigate to server directory:
```bash
cd server
```

2. Install dependencies:
```bash
npm install
```

3. Create `.env` file from `.env.example`:
```bash
cp .env.example .env
```

4. Configure `.env` file:
```env
SUPABASE_URL=your_supabase_url_here
SUPABASE_KEY=your_supabase_anon_key_here
ENCRYPTION_KEY=your-secret-encryption-key-32ch
PORT=3000
HOST=0.0.0.0
```

5. Start the server:
```bash
npm start
```

## Client Setup (Windows)

### Prerequisites
- Go 1.21 or higher
- Windows operating system

### Installation

1. Navigate to client directory:
```bash
cd client
```

2. Install dependencies:
```bash
go mod download
```

3. Build the client:
```bash
go build -o sqlbots-client.exe main.go
```

4. Run the client:
```bash
# Using environment variable
set API_KEY=your_api_key_here
sqlbots-client.exe

# Or using command line argument
sqlbots-client.exe your_api_key_here
```

### Environment Variables (Optional)

- `API_KEY`: Your API key (required if not provided as argument)
- `SERVER_URL`: Server URL (default: http://localhost:3000)
- `ENCRYPTION_KEY`: Encryption key (must match server)

## Database Schema

### Users Table
- `id`: UUID (Primary Key)
- `api_key`: Text (Unique)
- `license_id`: UUID (Foreign Key to licenses table)

### Licenses Table
- `id`: UUID (Primary Key)
- `expires_at`: Timestamp
- `user_id`: UUID (Foreign Key to users table)

### Machines Table
- `id`: UUID (Primary Key)
- `name`: Text (Machine name)
- `machine`: Text (Machine identifier)
- `api_key`: Text (Associated API key)
- `ram`: Integer (RAM in GB)
- `cores`: Integer (CPU cores)
- `created_at`: Timestamp
- `updated_at`: Timestamp

## API Endpoints

### POST /heartbeat
Sends encrypted heartbeat data to the server.

**Request Body:**
```json
{
  "data": "encrypted_string"
}
```

**Response:**
```json
{
  "data": "encrypted_response_string"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00.000Z"
}
```

## Error Codes

- `LICENSE_EXPIRED`: User's license has expired
- `MACHINE_LIMIT_EXCEEDED`: User has reached maximum machine limit (3)
- `INVALID_API_KEY`: Provided API key is invalid

## License

ISC

# sqlbot-client-server

