# Mentorship Backend

A Go-based backend for a mentorship platform.

## Deployment

### Vercel

1. Create a new project on Vercel
2. Connect your GitHub repository
3. Configure environment variables in Vercel dashboard:
   - DATABASE_URL
   - JWT_SECRET
   - PORT (optional, defaults to 3000)

### Railway

1. Create a new project on Railway
2. Connect your GitHub repository
3. Configure environment variables in Railway dashboard:
   - DATABASE_URL
   - JWT_SECRET
   - PORT (optional, defaults to 3000)

### Docker

1. Build the Docker image:
   ```bash
   docker build -t mentorship-backend .
   ```

2. Run the container:
   ```bash
   docker run -d -p 8080:8080 -e DATABASE_URL=your_db_url -e JWT_SECRET=your_secret mentorship-backend
   ```

## Environment Variables

- DATABASE_URL: PostgreSQL connection string
- JWT_SECRET: Secret key for JWT token signing
- PORT: Port number for the server (optional, defaults to 8080)

## Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

The server will run on http://localhost:8080 by default.
