FROM golang:1.14 as builder
WORKDIR /go/src/remember
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o remember


FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/remember/remember /remember
EXPOSE 1323
CMD ["/remember"]