#!/bin/bash

# Investify Development Script

echo "🚀 Starting Investify development environment..."

# Check if .env file exists and load it
if [ -f .env ]; then
  echo "🔑 Loading environment variables from .env file..."
  export $(grep -v '^#' .env | xargs)
fi

# Set Python path if virtual environment exists
if [ -d ".venv" ]; then
  echo "🐍 Using Python virtual environment: .venv"
  export PYTHONPATH=$(pwd)
  if [ "$(uname)" == "Darwin" ]; then
    # macOS
    source .venv/bin/activate
  else
    # Linux/Windows (Git Bash)
    source .venv/Scripts/activate 2>/dev/null || source .venv/bin/activate
  fi
fi

# Build the React frontend
cd frontend || exit
echo "📦 Building React frontend..."
npm run build
cd ..

# Build and run the Go backend
cd cmd || exit
echo "🔧 Building Go backend..."
go build -o investify
echo "🚀 Starting server..."
./investify
