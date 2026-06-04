#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "==================================="
echo "   Starting Northstar CMDB...      "
echo "==================================="

# Function to handle cleanup on exit (Ctrl+C)
cleanup() {
    echo ""
    echo "Stopping servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit
}

# Register the cleanup function for interrupt signals
trap cleanup SIGINT SIGTERM

# 0. Clean and prepare Frontend virtual routing caches and production builds
echo "-> Clearing Nuxt virtual route caches and production builds..."
rm -rf frontend/.nuxt
rm -rf frontend/.output

echo "-> Preparing Nuxt build directories..."
cd frontend
npm run postinstall
cd ..

# 1. Start the Backend
echo "-> Starting Go Backend on port 8000..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

# 2. Start the Frontend
echo "-> Starting NuxtJS Frontend on port 3000..."
cd frontend
# Using npm to start the dev server in the background
npm run dev &
FRONTEND_PID=$!
cd ..

echo "==================================="
echo " Both servers are booting up!"
echo " 🌐 Frontend : http://localhost:3000"
echo " ⚙️  Backend API: http://localhost:8000/docs"
echo "==================================="
echo "Press Ctrl+C to stop both servers."

# Wait indefinitely until interrupted
wait
