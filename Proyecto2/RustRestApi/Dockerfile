FROM rust:1.83 as builder

WORKDIR /usr/src/redis_client

COPY . .

RUN cargo install --path .

FROM alpine:3.14

COPY --from=builder /usr/local/cargo/bin/redis_client /usr/local/bin/redis_client

CMD ["redis_client"]