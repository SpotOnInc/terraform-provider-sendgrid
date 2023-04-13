# Use golang version 1.18 as the base image
FROM golang:1.18

# Set the repo name as a build argument (default to "app" if not provided)
ARG REPO_NAME=app

# Set the working directory to /app
WORKDIR /$REPO_NAME

# Copy the source code into the container at /app
COPY . .
