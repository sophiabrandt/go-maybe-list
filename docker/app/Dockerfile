# Build environment
# -----------------
FROM golang:1.15-alpine as build-env

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev git

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/web ./cmd/web
RUN go build -ldflags '-w -s' -a -o ./bin/admin ./cmd/admin

# Deployment environment
# ----------------------
FROM alpine:3.12

RUN apk update && apk add --no-cache bash

RUN adduser -D user
USER user

# Static Files
COPY ./ui /app/ui

# Binaries
COPY --from=build-env /app/bin/web /app/
COPY --from=build-env /app/bin/admin /app/

EXPOSE 8000
EXPOSE 6060

CMD [ "/app/web"]
