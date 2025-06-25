FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api

FROM gcr.io/distroless/base-debian12
COPY --from=build /go-api /go-api
EXPOSE 8080
ENTRYPOINT ["/go-api"]
