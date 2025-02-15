# Apptimus - Full-Stack Posts Application

## Project Overview
Apptimus is a comprehensive full-stack posts application featuring user authentication, post management, and static site generation.

## System Requirements
- Docker (latest version)
- Docker Compose (latest version)
- Git
- Supported Platforms:
  - macOS (Intel/Apple Silicon)
  - Linux (Ubuntu 20.04+, Fedora, CentOS)
  - Windows 10/11 (with WSL2)

## Prerequisites

### Docker Installation
- **macOS**: 
  ```bash
  brew install --cask docker


Ubuntu:
bashCopysudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io

Windows:
Install Windows Subsystem for Linux (WSL2)
Download and install Docker Desktop for Windows

Project Setup
1. Clone the Repository
cd apptimus
2. Environment Configuration
Create a .env file in the project root:
CopyDB_HOST=mysql_db
DB_USER=apptimus
DB_PASSWORD=1q2w3e
DB_NAME=apptimus_db
JWT_SECRET=your_secure_jwt_secret
NEXTAUTH_SECRET=your_secure_nextauth_secret
3. Initialize the Project
bashCopydocker-compose -f docker-compose.prod.yml up -d --build
4. Sample Data Loading
bashCopy# Connect to MySQL container
docker exec -it mysql_db mysql -u apptimus -p1q2w3e apptimus_db

# Insert sample data
INSERT INTO users (username, email, password_hash, token) VALUES 
('admin', 'admin@example.com', 'hashed_password', 'sample_token');

INSERT INTO posts (title, body, author, author_id) VALUES 
('Welcome to Apptimus', '<p>This is our first blog post!</p>', 'admin', 1);
5. Accessing the Application

Frontend: http://localhost:3000
Backend API: http://localhost:8080
Static Site: http://localhost:8081

Running Tests
bashCopy# Backend tests
docker-compose exec backend go test ./...

# Frontend tests
docker-compose exec frontend npm test
Stopping the Application
bashCopydocker-compose -f docker-compose.prod.yml down
Troubleshooting

Ensure Docker is running
Check that ports 3000, 8080, 8081 are free
Verify .env file credentials

Development Tools

Postman Collection: Use Apptimus_API_Collection.json for API testing
Static Site Generation: Run ./static-site/generate-static-site.sh

Contributing

Fork the repository
Create a feature branch
Commit changes
Push to the branch
Create a Pull Request