FROM golang:1.16-alpine as builder

RUN apk add --no-cache make

COPY ./ /workspace
WORKDIR /workspace

RUN make -C local archiver-local

FROM alpine as app
COPY --from=builder /workspace/local/archiver-local ./archiver-local
CMD ["./archiver-local"]
