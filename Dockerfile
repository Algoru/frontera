FROM golang:1.15 AS build

WORKDIR /go/src/github.com/Algoru/frontera
COPY . .

RUN go mod tidy
RUN go build

FROM alpine AS productive

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/Algoru/frontera .

CMD ["./frontera"] 