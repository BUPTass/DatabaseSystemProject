FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app .

FROM python:3.11

RUN pip install matplotlib networkx

COPY --from=builder /build/app /
COPY run.py /
RUN mkdir upload


ENTRYPOINT ["/app"]