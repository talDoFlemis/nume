FROM node:lts AS frontend_base

FROM frontend_base AS frontend_builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app
COPY frontend/ /app

RUN --mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/pnpm/store,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache \
	gcc \
	musl-dev \
	ca-certificates

COPY go.mod go.sum /app/

RUN --mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/go/pkg/mod/,target=/go/pkg/mod/ \
	go mod download -x

COPY . .
COPY --from=frontend_builder /app/dist /app/frontend/dist

RUN --mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/go/pkg/mod/,target=/go/pkg/mod/ \
	--mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/root/.cache/go-build,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -o server  -ldflags '-s -w -extldflags "-static"' ./cmd/web

FROM ubuntu:oracular AS user
RUN useradd -u 10001 scratchuser

FROM scratch
WORKDIR /app


COPY --from=builder /app/server ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=user /etc/passwd /etc/passwd

USER scratchuser
STOPSIGNAL SIGINT
EXPOSE 8080

ENV NUME_APP_ENVIRONMENT="prod"

CMD ["/app/server"]

