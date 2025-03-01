FROM golang:1.24 AS deps
WORKDIR /project
COPY api/go.mod api/go.sum api/
COPY fiber-tools/go.mod fiber-tools/go.sum fiber-tools/
COPY cmd/go.mod cmd/go.sum cmd/
COPY go.work .
RUN go mod download


FROM deps AS build
COPY api/ api/
COPY cmd/ cmd/
COPY fiber-tools/ fiber-tools/
RUN go build ./cmd/run-server


FROM ubuntu:noble
SHELL [ "/bin/bash", "-c" ]
WORKDIR /app
COPY --from=build /project/run-server .