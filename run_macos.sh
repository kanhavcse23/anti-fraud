#!/bin/bash

# Create and activate Python virtual environment
python3 -m venv env
source env/bin/activate

# Update system and install required dependencies
brew update

# Install Docker Desktop (Docker for macOS includes Docker Compose)
brew install --cask docker

# Ensure Docker is running
open /Applications/Docker.app
echo "Waiting for Docker to start..."
until docker info > /dev/null 2>&1; do
    sleep 1
done

# Install Go
brew install go

# Ensure dependencies are up to date
go mod tidy

# Stop and clean up any existing Docker containers
docker-compose down
docker system prune -a -f

# Build and start Docker containers
docker-compose up --build

# Deactivate Python virtual environment
deactivate