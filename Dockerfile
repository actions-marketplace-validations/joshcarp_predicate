FROM golang:buster
WORKDIR /usr/app
ADD . .
RUN go build -o gh-issue-automation .
RUN mv gh-issue-automation /bin/gh-issue-automation
RUN apt-get update && apt-get install ca-certificates && update-ca-certificates
ENTRYPOINT gh-issue-automation