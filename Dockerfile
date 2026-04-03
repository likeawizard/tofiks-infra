# Stage 1: Build tofiks binary
FROM golang:1.26 AS builder

WORKDIR /build
RUN git clone https://github.com/likeawizard/tofiks.git .
ARG TOFIKS_REF=master
RUN git checkout ${TOFIKS_REF}
RUN GOAMD64=v3 CGO_ENABLED=0 go build -gcflags=-B -o /tofiks cmd/tofiks/main.go

# Stage 2: Runtime with lichess-bot + tofiks
FROM python:3.12-slim

WORKDIR /app

# Install lichess-bot dependencies
COPY lichess-bot/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy lichess-bot source
COPY lichess-bot/ .

# Copy tofiks binary and opening book
COPY --from=builder /tofiks engines/tofiks
COPY engines/tofiks.bin engines/tofiks.bin

# Copy config template (token injected at runtime via entrypoint)
COPY config.yml config.yml.template

# Entrypoint: substitute LICHESS_TOKEN into config, then run the bot
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
