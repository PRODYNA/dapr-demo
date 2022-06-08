FROM golang AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o goapp

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=builder /app/goapp /goapp
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["/goapp"]
