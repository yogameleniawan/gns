FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache git gettext ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP_PORT
ARG JWT_SECRET
ARG JWT_REFRESH_SECRET
ARG JWT_ACCESS_TOKEN_EXPIRATION
ARG JWT_REFRESH_TOKEN_EXPIRATION
ARG CACHE_HOST
ARG CACHE_PORT
ARG CACHE_DB
ARG SUPABASE_URL
ARG SUPABASE_KEY

ENV APP_PORT=${APP_PORT}
ENV JWT_SECRET=${JWT_SECRET}
ENV JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
ENV JWT_ACCESS_TOKEN_EXPIRATION=${JWT_ACCESS_TOKEN_EXPIRATION}
ENV JWT_REFRESH_TOKEN_EXPIRATION=${JWT_REFRESH_TOKEN_EXPIRATION}
ENV CACHE_HOST=${CACHE_HOST}
ENV CACHE_PORT=${CACHE_PORT}
ENV CACHE_DB=${CACHE_DB}
ENV SUPABASE_URL=${SUPABASE_URL}
ENV SUPABASE_KEY=${SUPABASE_KEY}

RUN envsubst < ./config/config.template.yaml > ./config/config.staging.yaml

RUN go build -o main cmd/api/main.go

# -------------------------------------------------
# Stage 2: Runtime (alpine)
# -------------------------------------------------
FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

CMD ["./main", "--env=staging"]