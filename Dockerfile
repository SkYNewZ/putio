FROM golang:1.14

# Install packr2
RUN go get -u github.com/gobuffalo/packr/v2/packr2

WORKDIR /go/src/github.com/SkYNewZ/putio
COPY go.* ./
RUN go mod download

COPY . .
# Embed static files
RUN PATH=$PATH:$(go env GOPATH)/bin packr2 --verbose

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /putio .


FROM scratch
LABEL description="Putio HTTP interface"
LABEL maintainer="Quentin Lemaire <quentin@lemairepro.fr>"

ENV PORT 8080
WORKDIR /app

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /putio /app/putio

ENTRYPOINT ["/app/putio"]
