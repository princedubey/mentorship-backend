# Mentorship Backend

A Go-based backend for a mentorship platform.

## Deployment

### Vercel

1. Create a new project on Vercel
2. Connect your GitHub repository
3. Add the following environment variables in Vercel:
   - `DATABASE_URL`: Your PostgreSQL database URL
   - `JWT_SECRET`: Your JWT secret key
   - `PORT`: Set to 3000 (Vercel default)
   - `CLOUDINARY_CLOUD_NAME`: Your Cloudinary cloud name
   - `CLOUDINARY_API_KEY`: Your Cloudinary API key
   - `CLOUDINARY_API_SECRET`: Your Cloudinary API secret
   - `FIREBASE_PROJECT_ID`: Your Firebase project ID
   - `FIREBASE_PRIVATE_KEY`: Your Firebase private key
   - `FIREBASE_CLIENT_EMAIL`: Your Firebase client email

4. Deploy the project

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
