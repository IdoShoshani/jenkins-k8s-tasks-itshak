# Docker Multi-App Demo Project

This project contains three simple web applications built with different technologies, each containerized using Docker.

## Project Structure

```
jenkins-k8s-tasks-itshak/
├── python-app/
│   ├── app.py
│   ├── requirements.txt
│   └── Dockerfile
├── node-app/
│   ├── app.js
│   ├── package.json
│   └── Dockerfile
└── go-app/
    ├── main.go
    └── Dockerfile
```

## Applications Overview

### Python Flask App

- Simple Flask web server
- Runs on port 5000
- Returns a "Hello from Python Flask!" message

### Node.js Express App

- Express.js web server
- Runs on port 3000
- Returns a "Hello from Node.js Express!" message

### Go App

- Basic Go web server
- Runs on port 8080
- Returns a "Hello from Go!" message

## Prerequisites

- Docker installed on your machine
- Git (optional, for cloning the repository)

## Building and Running the Applications

### Python App

```bash
cd jenkins-k8s-tasks-itshak//python-app
docker build -t python-app .
docker run -p 5000:5000 python-app
```

Access at: http://localhost:5000

### Node.js App

```bash
cd jenkins-k8s-tasks-itshak//node-app
docker build -t node-app .
docker run -p 3000:3000 node-app
```

Access at: http://localhost:3000

### Go App

```bash
cd jenkins-k8s-tasks-itshak//go-app
docker build -t go-app .
docker run -p 8080:8080 go-app
```

Access at: http://localhost:8080

## Docker Commands Reference

Stop a container:

```bash
docker stop <container_id>
```

List running containers:

```bash
docker ps
```

List all containers (including stopped):

```bash
docker ps -a
```

Remove a container:

```bash
docker rm <container_id>
```

List Docker images:

```bash
docker images
```

Remove a Docker image:

```bash
docker rmi <image_name>
```

## Troubleshooting

1. Port already in use:

   - Stop any running containers using the same port
   - Or change the port mapping in the docker run command (e.g., -p 5001:5000)

2. Build fails:

   - Check if Docker daemon is running
   - Ensure all required files are present in the correct directories
   - Check internet connectivity for downloading dependencies

3. Container exits immediately:
   - Check container logs: `docker logs <container_id>`
   - Ensure the application is properly configured to run in the foreground

## Contributing

Feel free to submit issues and enhancement requests!
