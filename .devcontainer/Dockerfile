# Global ARGs
ARG GO_VERSION=1.23
ARG DISTRIB=alpine

FROM golang:${GO_VERSION}-${DISTRIB}

# Install required packages for Go development with CGO
RUN apk update && apk add --no-cache \
    git \
    gcc \
    musl-dev \
    bash \
    bash-completion

# Go setup
ENV CGO_ENABLED=1

# Set working directory
WORKDIR /workspace

# Copy only go.mod first
COPY go.mod .

# Download dependencies (this will create go.sum)
RUN go mod download

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Bash completion setup
RUN echo "source /usr/share/bash-completion/bash_completion" >> /etc/bash.bashrc

# Copy the rest of the source code
COPY . .
