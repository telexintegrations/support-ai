# Use an official Python image
FROM python:3.10-slim

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Install MongoDB
RUN apt-get update && apt-get install -y mongodb && rm -rf /var/lib/apt/lists/*

# Create MongoDB data directory
RUN mkdir -p /data/db

# Install ChromaDB
RUN pip install chromadb==0.5.5

# Expose MongoDB and ChromaDB ports
EXPOSE 27017 8000

# Start both MongoDB and ChromaDB
CMD mongod --fork --logpath /var/log/mongodb.log && chromadb run --host 0.0.0.0 --port 8000
