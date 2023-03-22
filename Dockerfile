FROM golang:1.19-alpine as build-env
RUN apk --no-cache add build-base gcc
ADD . /src
RUN cd /src && go mod tidy && go build -o voomkit

FROM alpine
WORKDIR /app
RUN apk --no-cache add build-base gcc
COPY --from=build-env /src/voomkit /app/voomkit
CMD /app/voomkit
