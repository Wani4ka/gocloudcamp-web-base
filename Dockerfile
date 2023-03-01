## Build
FROM golang:1.19.6-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /gocloudcamp

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gocloudcamp /gocloudcamp

USER nonroot:nonroot

ENTRYPOINT ["/gocloudcamp"]
