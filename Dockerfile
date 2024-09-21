# Use the official MongoDB image as the base image
FROM mongo:latest

# Set environment variables (these can be overridden in docker-compose.yml)
ENV MONGO_INITDB_ROOT_USERNAME=admin
ENV MONGO_INITDB_ROOT_PASSWORD=password

# Expose the default MongoDB port
EXPOSE 27017

# The CMD is not necessary here as it's inherited from the base image
# but you can uncomment it if you want to be explicit
# CMD ["mongod"]