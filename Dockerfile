FROM golang:1.27.1-alpine3.23 AS build_image

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine:3.23

EXPOSE 6090

WORKDIR /project

COPY --from=build_image /src/app .
COPY --from=build_image /src/frontend ./frontend

ENTRYPOINT [ "./app" ]