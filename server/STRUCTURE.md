# Server Structure

## Project Organization

```
server/
├── index.js                 # Main entry point
├── config.js                # Configuration management
├── package.json             # Dependencies
├── db/
│   └── supabase.js         # Supabase client initialization
├── services/
│   ├── encryption.js        # Encryption/decryption service
│   ├── userService.js       # User-related operations
│   ├── machineService.js    # Machine-related operations
│   └── licenseService.js    # License validation operations
└── routes/
    ├── health.js           # Health check route
    └── heartbeat.js       # Heartbeat route handler
```

## Module Descriptions

### config.js
Centralized configuration management for:
- Supabase connection (URL and key)
- Encryption key
- Server settings (port and host)

### db/supabase.js
Initializes and exports the Supabase client instance.

### services/encryption.js
Handles all encryption and decryption operations:
- `decryptData()` - Decrypts incoming heartbeat data
- `encryptData()` - Encrypts response data

### services/userService.js
User-related database operations:
- `validateApiKey()` - Validates API key and retrieves user info

### services/licenseService.js
License validation operations:
- `isLicenseExpired()` - Checks if a license is expired
- `getLicenseExpiration()` - Retrieves license expiration date
- `validateLicense()` - Comprehensive license validation

### services/machineService.js
Machine-related database operations:
- `getMachineCount()` - Gets count of machines for an API key
- `getMachineByIdentifier()` - Retrieves machine by API key and machine ID
- `registerMachine()` - Registers a new machine
- `updateMachine()` - Updates machine hardware information
- `verifyHardware()` - Verifies if hardware matches stored data

### routes/health.js
Health check endpoint:
- `GET /health` - Returns server status

### routes/heartbeat.js
Heartbeat endpoint handler:
- `POST /heartbeat` - Processes encrypted heartbeat data
  - Validates API key
  - Checks license expiration
  - Manages machine registration
  - Verifies hardware information
  - Enforces machine limit (max 3)

### index.js
Main application entry point:
- Initializes Fastify server
- Registers all routes
- Starts the server

## Benefits of Modular Structure

1. **Separation of Concerns**: Each module has a single responsibility
2. **Maintainability**: Easier to locate and modify specific functionality
3. **Testability**: Individual modules can be tested in isolation
4. **Reusability**: Services can be reused across different routes
5. **Scalability**: Easy to add new routes and services

