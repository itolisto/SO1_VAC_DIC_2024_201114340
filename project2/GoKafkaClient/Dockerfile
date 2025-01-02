FROM golang:alpine AS build

WORKDIR /usr/src/app


COPY go.mod go.sum ./

# only redownloading dependencies them in subsequent builds if they change
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 9092

CMD ["app"]