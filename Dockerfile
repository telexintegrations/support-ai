
FROM python:3.10-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN pip install --no-cache-dir chromadb==0.5.5

EXPOSE 8000

CMD ["chromadb", "run", "--host", "0.0.0.0", "--port", "8000"]
