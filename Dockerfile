## Build
FROM golang:1.19.6-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /gocloudcamp

## Deploy
FROM scratch

WORKDIR /

COPY --from=build /gocloudcamp /gocloudcamp

ENTRYPOINT ["/gocloudcamp"]
