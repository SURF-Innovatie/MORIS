# Builder stage
FROM golang:1.25.4-alpine AS backend-gen

RUN apk add --no-cache git

# Install swag tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app/apps/backend

# Copy backend go modules
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download

# Copy backend source
COPY apps/backend .

# Run code generation
# 1. Generate swagger docs (ent needs this)
RUN swag init -g cmd/dev/main.go --output api/swag-docs --parseDependency --parseInternal
# 2. Generate ent code
RUN go generate ./ent/...
# 3. Generate events (if needed, though frontend mostly needs swagger)
RUN go generate ./internal/domain/events/...


# Frontend Builder stage
FROM node:20-alpine AS builder

# Enable pnpm
RUN corepack enable

WORKDIR /app

# Copy root workspace configuration
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json tsconfig.base.json ./

COPY apps/frontend/package.json ./apps/frontend/package.json
RUN pnpm install --frozen-lockfile --filter frontend...

# Copy source code
COPY apps/frontend ./apps/frontend

# Copy generated backend artifacts to where frontend script expects them
# Script looks at ../../backend/api/swag-docs/swagger.json 
# relative to apps/frontend/scripts
# So in container: /app/apps/backend/api/swag-docs/swagger.json
COPY --from=backend-gen /app/apps/backend/api/swag-docs /app/apps/backend/api/swag-docs

# Build the application
WORKDIR /app/apps/frontend
RUN pnpm run api:generate
RUN pnpm run build

# Final stage
FROM nginx:alpine

# Copy built assets from builder
COPY --from=builder /app/apps/frontend/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
