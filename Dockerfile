FROM golang:1.14
WORKDIR /go/src/github.com/SkYNewZ/putio

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /putio .


FROM scratch
LABEL description="Putio HTTP interface"
LABEL maintainer="Quentin Lemaire <quentin@lemairepro.fr>"

ENV PORT 8080
WORKDIR /app

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /putio /app/putio
COPY --from=0 /go/src/github.com/SkYNewZ/putio/templates /app/templates
COPY --from=0 /go/src/github.com/SkYNewZ/putio/assets /app/assets

ENTRYPOINT ["/app/putio"]
