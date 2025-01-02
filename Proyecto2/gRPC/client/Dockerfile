FROM golang:alpine AS build

WORKDIR /usr/src/app/client

COPY . .

COPY ../ProtoBuffer ../

COPY ../go.mod ../go.sum ../

RUN cd ..

# only redownloading dependencies them in subsequent builds if they change
RUN go mod download && go mod verify

RUN go build -v -o /usr/local/bin/app ./client/...

EXPOSE 8000

CMD ["app"]