FROM golang:1.23-alpine AS gobuild

RUN apk --no-cache upgrade \
  && apk --no-cache add \
    git \
    ca-certificates \
;

ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

RUN go mod download

COPY cmd /app/cmd
COPY pkg /app/pkg

RUN go build  -o /app/build/server ./cmd/api

FROM scratch

EXPOSE 80

COPY --from=gobuild /app/build/server ./
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY config /config

CMD [ "/server" ]

