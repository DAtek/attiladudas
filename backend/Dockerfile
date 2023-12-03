FROM golang:1.21 as deps
WORKDIR /project
COPY api/go.mod api/go.sum api/
COPY db/go.mod db/go.sum db/
COPY fiber-tools/go.mod fiber-tools/go.sum fiber-tools/
COPY cmd/go.mod cmd/go.sum cmd/
COPY go.work .
RUN go mod download


FROM deps as build
COPY api/ api/
COPY cmd/ cmd/
COPY db/ db/
COPY fiber-tools/ fiber-tools/
RUN go build ./cmd/run-server


FROM ubuntu:jammy
SHELL [ "/bin/bash", "-c" ]
# must match with nginx's user, see https://hub.docker.com/_/nginx/
RUN addgroup --gid 101 backendgroup
RUN useradd -u 101 -g backendgroup backend
WORKDIR /app
COPY --from=build /project/run-server .