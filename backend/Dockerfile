FROM golang:latest as deps

ENV GOPATH=
COPY go.mod go.sum ./
RUN go mod download


FROM deps as build

WORKDIR /build
COPY . .
RUN go build ./cmd/server


FROM ubuntu:latest
# must match with nginx's user, see https://hub.docker.com/_/nginx/
RUN addgroup --gid 101 backendgroup
RUN useradd -u 101 -g backendgroup backend
WORKDIR /app
COPY --from=build /build/server .
CMD [ "./server" ]