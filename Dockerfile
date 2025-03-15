# Use MongoDB 7.0 official image
FROM mongo:7.0

# Set up working directory inside the container
WORKDIR /data/db

# Expose MongoDB's default port
EXPOSE 27017

# Start MongoDB when the container runs
CMD ["mongod", "--bind_ip_all"]
