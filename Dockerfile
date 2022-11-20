FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY ./ ./

RUN apk add --no-cache --update make && make build

# ---

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/dist/app /

EXPOSE 8080
ENTRYPOINT ["/app"]
