# FROM golang:1.16-alpine AS build_base

# RUN apk add --no-cache git

# WORKDIR /tmp/library-app

# COPY go.mod .
# COPY go.sum .

# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 go test -v

# RUN go build -o ./out/library-app .

# FROM alpine:3.9 
# RUN apk add ca-certificates

# COPY --from=build_base /tmp/library-app/out/library-app /app/library-app

# EXPOSE 8080

# CMD ["/app/library-app"]

FROM golang:1.16-alpine AS build

WORKDIR /src/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go test -v
RUN go build -o /athenaeum/athenaeumserver .

FROM alpine

RUN mkdir /resources
RUN mkdir /libhome
RUN mkdir /database

COPY resources/ /resources/
COPY --from=build /athenaeum/athenaeumserver /athenaeum/athenaeumserver

ENV LOG_LEVEL="info"

ENV STATIC_DIR="/resources/static/"
ENV TEMPLATE_DIR="/resources/templates/"

ENV DB_LOCATION="/database/"
ENV LIBRARY_HOME="/libhome/"

ENV ADDRESS=":8080"
EXPOSE 8080

ENTRYPOINT ["/athenaeum/athenaeumserver"]
CMD [""]
