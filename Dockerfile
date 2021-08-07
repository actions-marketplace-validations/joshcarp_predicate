FROM golang:buster
WORKDIR /usr/app
ADD . .
RUN go build -o predicate .
RUN mv predicate /bin/predicate
RUN apt-get update
RUN apt-get install ca-certificates
RUN go get github.com/fullstorydev/grpcurl/cmd/grpcurl
RUN update-ca-certificates
ENTRYPOINT predicate