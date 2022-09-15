##
## BUILD
##
FROM golang:1.18-alpine3.15 AS build

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN chmod +x build.sh

RUN ./build.sh

##
## Deploy
##
FROM alpine:3.15

RUN adduser -D nonroot

WORKDIR /

COPY --from=build /app/bin/school-book-storage-service /school-book-storage-service

EXPOSE 9090

USER nonroot

ENTRYPOINT [ "/school-book-storage-service" ]