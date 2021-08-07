FROM golang:buster
WORKDIR /usr/app
ADD . .
RUN go build -o predicate .
RUN mv predicate /bin/predicate
RUN apt-get update && apt-get install ca-certificates && update-ca-certificates
ENTRYPOINT predicate