FROM 'golang:1.13-alpine'

RUN adduser -D -u 1000 -g heroes_api heroes_api

RUN mkdir /app

WORKDIR /app

RUN apk add git
