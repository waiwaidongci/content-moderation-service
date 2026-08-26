FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /content-moderation-service ./cmd/content-moderation-service
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /content-moderation-service /content-moderation-service
EXPOSE 8083
ENTRYPOINT ["/content-moderation-service"]
