FROM golang:1.20 AS builder
WORKDIR /pr1/
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./app ./main.go

FROM alpine:latest
WORKDIR /pr1
COPY --from=builder /pr1/.env .
COPY --from=builder /pr1/app .
ENTRYPOINT [ "/pr1/app" ]