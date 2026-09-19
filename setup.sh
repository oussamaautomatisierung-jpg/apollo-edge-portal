#!/bin/bash

echo "Apollo Edge Portal - Setup"
echo "----------------------------"

# Check that Docker is installed
if ! command -v docker &> /dev/null; then
	echo "Error: Docker is not installed or not in PATH."
	echo "Please install Docker Desktop and try again."
	exit 1
fi

# Check that the Docker daemon is running
if ! docker info &> /dev/null; then
	echo "Error: Docker is installed but not running."
	echo "Please start Docker Desktop and try again."
	exit 1
fi

echo "Docker is running."
echo "Starting the stack (mosquitto + go-portal)..."
echo ""

docker compose up --build
