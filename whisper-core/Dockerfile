
## Build golang
FROM golang:1.21.1 AS build_go

WORKDIR /app
RUN cd /app
COPY ./ ./
RUN go mod download
RUN go build -o /whispercore ./cmd/whisper-core/main.go

## DEPLOYMENT
FROM golang:1.21.1

WORKDIR /app
COPY --from=build_go /whispercore /app/whispercore
EXPOSE 8080
CMD [ "/app/whispercore" ]