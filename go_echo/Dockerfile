FROM golang:1.18-alpine

RUN apk update && apk add git curl

WORKDIR /usr/src/app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download && go mod verify

COPY . .

#
# Development
#
RUN go install -v -n -a github.com/cosmtrek/air@v1.29.0
CMD ["./bin/air", "-c", ".air.toml"]


#
# Production
#
#RUN go build -v -o /usr/local/bin/app
#CMD ["/usr/local/bin/app"]