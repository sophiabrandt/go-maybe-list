# Build environment
# -----------------
FROM golang:1.15-alpine as build-env

WORKDIR /app

RUN apk update && apk add --no-cache sqlite gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o ./bin/web ./cmd/web
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o ./bin/admin ./cmd/admin

# Deployment environment
# ----------------------
FROM alpine:3.12

WORKDIR /app

RUN apk update && apk add --no-cache bash

RUN adduser -D user && chown -R user:user /app
USER user

# Static Files & database
COPY --chown=user:user ./ui ./ui

# Binaries
COPY --from=build-env --chown=user:user /app/bin/web .
COPY --from=build-env --chown=user:user /app/bin/admin .

EXPOSE 4000

CMD [ "./web", "-secret=$COOKIE_SECRET" ]
