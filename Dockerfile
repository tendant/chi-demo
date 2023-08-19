FROM golang:1.20.5-alpine AS build

RUN apk update && \
    apk add --update openntpd && \
    ntpd && \
    apk upgrade && \
    apk add --no-cache alpine-sdk git make

WORKDIR /app

# Cache go mod dependencies
COPY go.mod ./
RUN go mod download

COPY . .

RUN make

FROM golang:1.20.5-alpine

WORKDIR /app

COPY --from=build /app/dist/ /app/

RUN chmod a+x /app/*

# CMD ["/app/main"]
