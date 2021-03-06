# build on heroku
FROM golang:1.10.4 AS build-env

ADD . /src
RUN cd /src && dep ensure && GOOS=linux GOARCH=386 go build -o goapp

FROM alpine
ENV PORT=$PORT
WORKDIR /app
COPY --from=build-env /src/goapp /app/
CMD /app/goapp