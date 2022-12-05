FROM golang:1.18 as builder
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ARG version=dev

RUN mv bin/migrate.linux-amd64 migrate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" -o /bin/movierama ./cmd/movierama/main.go
RUN adduser -u 5003 --gecos '' --disabled-password --no-create-home movierama

FROM alpine:3.17.0
RUN apk --no-cache add tzdata
COPY --from=builder /bin/movierama .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /src/migrate .
COPY --from=builder /src/migration ./migration
COPY --from=builder /src/script/migrate_and_run.sh .

USER movierama
ENTRYPOINT ["./migrate_and_run.sh"]
