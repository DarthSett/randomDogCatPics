FROM golang
MAINTAINER sourav241196@gmail.com
WORKDIR /app
ENV GO111MODULE=on
COPY ./go.mod ./go.sum ./
RUN go mod tidy
RUN go mod download
RUN go mod verify
COPY ./ .
RUN CGO_ENABLED=0 go build -v -o ./bin/app .
EXPOSE 3000

FROM alpine:latest
COPY --from=0 /app/bin/app .
CMD ["./app"]