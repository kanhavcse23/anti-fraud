#!/bin/bash
python3 -m venv env
source env/bin/activate

# sudo apt update
# sudo apt install -y ca-certificates curl gnupg
# sudo mkdir -p /etc/apt/keyrings
# curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
# echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
# sudo apt update
# sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
# sudo apt install -y golang-go
go mod tidy
docker-compose down

# docker system prune -a -f

sudo systemctl start docker
# sudo apt install -y docker-compose
docker-compose up --build
deactivate