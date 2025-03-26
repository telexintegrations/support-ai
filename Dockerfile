FROM python:3.9-slim

ENV DEBIAN_FRONTEND=noninteractive

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Install ChromaDB
RUN pip install --no-cache-dir chromadb==0.5.5 onnxruntime protobuf

# Expose ChromaDB API port
EXPOSE 8000

# Start ChromaDB server
CMD ["chroma", "run", "--host", "0.0.0.0", "--port", "8000"]
