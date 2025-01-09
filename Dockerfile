# Fetch dependencies
FROM golang:1.23.1 AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

# Generate templ
FROM ghcr.io/a-h/templ:latest AS templ-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build CSS
FROM node:20-alpine AS css-stage
WORKDIR /app
COPY --from=templ-stage /app .
RUN npm install
RUN npx tailwindcss -i ./view/index.css -o ./public/styles.css --minify

# Build Go binary
FROM golang:1.23.1 AS build-stage
WORKDIR /app
COPY --from=fetch-stage /go/pkg/mod /go/pkg/mod
COPY --from=templ-stage /app .
COPY --from=css-stage /app/public/styles.css ./public/styles.css
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main

# Final minimal image
FROM alpine
WORKDIR /app
COPY --from=build-stage /app/main .
COPY --from=css-stage /app/public ./public
COPY --from=build-stage /app/.env .
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT ["/app/main"]
