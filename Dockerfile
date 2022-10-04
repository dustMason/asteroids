FROM golang:1.19-buster

ADD . /src
WORKDIR /src
RUN GOOS=linux go build -o /bin/asteroids.

ENTRYPOINT ["/bin/asteroids"]
