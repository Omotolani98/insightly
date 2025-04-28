#!/bin/bash

# Insightly Starter Script

echo "🚀 Starting Redis and Postgres..."
make up

echo "⏳ Waiting for containers to initialize..."
sleep 5

echo "🛠 Starting Ingest server..."
make ingest &

echo "🛠 Starting Summarizer worker..."
make summarizer &

echo "🛠 Starting Query server..."
make query &

echo "🛠 Starting API server..."
make api &

echo "🎉 All Insightly services are now starting!"