FROM golang:1.27.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /task-manager-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /task-manager-api /task-manager-api
EXPOSE 8080
ENTRYPOINT ["/task-manager-api"]
