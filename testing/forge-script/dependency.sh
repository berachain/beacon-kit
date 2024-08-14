#!/bin/sh
apk update && apk add --no-cache nodejs npm

npm --version
npm install -g bun

bun --version

cd /app/contracts && bun install

echo "Bun installation complete!"
