FROM golang:1.19-alpine AS builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o dist/app cmd/server/main.go

# ---

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/dist/app /

EXPOSE 8080
ENTRYPOINT ["/app"]
