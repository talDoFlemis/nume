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

RUN --mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/go/pkg/mod/,target=/go/pkg/mod/ \
  --mount=type=cache,id=s/78f32e04-4dcf-4cb1-9282-97744b29639c-/root/.cache/go-build,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build -o server  -ldflags '-s -w -extldflags "-static"' ./cmd/ssh

FROM ubuntu:oracular AS user
RUN useradd -u 10001 scratchuser && mkdir -p /home/scratchuser/.ssh && chown scratchuser:scratchuser /home/scratchuser/.ssh

FROM scratch
WORKDIR /app


COPY --from=builder /app/server ./
COPY .ssh /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=user /etc/passwd /etc/passwd
COPY --from=user --chown=scratchuser:scratchuser /home/scratchuser/.ssh /app/.ssh

USER scratchuser
STOPSIGNAL SIGINT
EXPOSE 8888

ENV NUME_APP_ENVIRONMENT="prod"

CMD ["/app/server"]

