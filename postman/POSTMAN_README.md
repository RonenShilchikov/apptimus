# Apptimus API Postman Collection

## Setup Instructions

1. Download Postman from https://www.postman.com/downloads/
2. Import the `Apptimus_API_Collection.json` file
3. Set up environment variables:
   - `baseUrl`: Your API base URL (default: http://localhost:8080)
   - `token`: JWT token received after login
   - `userId`: User ID for deletion
   - `postId`: Post ID for editing/deletion

## Authentication Flow
1. Register a new user via the Register endpoint
2. Login to receive a JWT token
3. Copy the token and set it in the `token` variable
4. Use the token for authenticated requests

## Common Troubleshooting
- Ensure the backend is running
- Check network configuration
- Verify JWT token is valid and not expired