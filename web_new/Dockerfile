# Stage 1: Build the SvelteKit application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package definitions and install dependencies
COPY package.json package-lock.json* ./
RUN npm install

# Copy the rest of the application source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Create the production image
FROM node:20-alpine

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/.svelte-kit/output ./build

# Copy production dependencies
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Set environment variables for the Node.js server
ENV HOST=0.0.0.0
ENV PORT=3000

# Expose the port the application will run on
EXPOSE 3000

# The command to start the Node.js server
# The entry point is index.js inside the 'build' directory
CMD ["node", "build/index.js"]
