FROM golang:1.17.2
LABEL maintainer="agaesser@gmail.com"

ENV TERM linux
RUN apk --no-cache add apache2-utils