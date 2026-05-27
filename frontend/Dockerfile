# Step 1: Build Angular application via Node environment
FROM node:20-alpine AS build

WORKDIR /app

# [CRITICAL]: Copy package files first to cache dependencies layers
COPY package*.json ./

# Download and install @angular/core and all dependencies inside Docker container
RUN npm install

# Copy the rest of the frontend application source code
COPY . .

# Compile code to production static files
RUN npm run build -- --configuration=production

# Step 2: Serve static files using lightweight Nginx web server
FROM nginx:alpine

# Copy compiled files from build stage to Nginx web directory
COPY --from=build /app/dist/example.com/browser /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
