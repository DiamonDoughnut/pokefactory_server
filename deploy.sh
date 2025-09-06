#!/bin/bash

# Docker Hub deployment script
DOCKER_USERNAME="yourusername"  # Replace with your Docker Hub username
IMAGE_NAME="pokefactory-server"
VERSION="v1.0.0"

echo "Building Docker image..."
docker build -t $DOCKER_USERNAME/$IMAGE_NAME:$VERSION .
docker build -t $DOCKER_USERNAME/$IMAGE_NAME:latest .

echo "Logging into Docker Hub..."
docker login

echo "Pushing to Docker Hub..."
docker push $DOCKER_USERNAME/$IMAGE_NAME:$VERSION
docker push $DOCKER_USERNAME/$IMAGE_NAME:latest

echo "Deployment complete!"
echo "Your partner can now run:"
echo "docker pull $DOCKER_USERNAME/$IMAGE_NAME:latest"