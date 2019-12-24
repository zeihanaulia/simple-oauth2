FROM golang:1.13-alpine AS build_base
RUN apk update && apk add --no-cache ca-certificates git make bash
WORKDIR /usr/src/app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/oauth2

FROM alpine
RUN apk add ca-certificates
COPY --from=server_builder /usr/src/app/bin/oauth2 .

# copy templates
COPY client/templates/* ./client/templates/
COPY services/authorization/templates/* ./services/authorization/templates/
COPY services/protected/templates/* ./services/protected/templates/

EXPOSE 8080
ENTRYPOINT ["./oauth2"]
CMD ["all"]