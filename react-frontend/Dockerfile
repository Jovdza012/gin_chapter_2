# Use an older Node.js version to avoid OpenSSL issues
FROM node:16

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json first
COPY package*.json ./

# Copy the rest of the application code
COPY . .

# Expose the port the React app will run on
EXPOSE 3000

# Start the React app
CMD ["npm", "start"]