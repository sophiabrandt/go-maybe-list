# Build environment
# -----------------
FROM golang:1.22-alpine as compile-image

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o ./bin/web ./cmd/web
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o ./bin/admin ./cmd/admin

# Deployment environment
# ----------------------
FROM alpine:3.20 as runtime-image

WORKDIR /app

RUN apk update && apk add --no-cache bash sqlite

RUN adduser -D user && chown -R user:user /app
USER user

# Static Files
COPY --chown=user:user ./ui ./ui

# Binaries
COPY --from=compile-image --chown=user:user /app/bin/web .
COPY --from=compile-image --chown=user:user /app/bin/admin .

EXPOSE 4000

RUN ./admin -action="migrate"

CMD [ "./web" ]
