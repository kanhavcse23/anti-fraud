#!/bin/bash
gh repo clone kanhavcse23/anti-fraud
# Activate the virtual environment
source rishav_project_env/bin/activate

# Your commands here
# Example:
sudo apt update
sudo apt install -y golang-go
go mod tidy
sudo systemctl start docker
docker-compose up --build

# Deactivate the virtual environment
# deactivate
