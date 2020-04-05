FROM golang:1.14 as builder
WORKDIR /go/src/remember
COPY go.* ./
RUN go mod download
COPY . ./
# -mod=readonly
RUN CGO_ENABLED=0 GOOS=linux go build  -v -o remember


FROM alpine:3
RUN apk --no-cache add ca-certificates
EXPOSE 1323
ENV PORT=1323
CMD ["/remember"]
COPY --from=builder /go/src/remember/remember /remember