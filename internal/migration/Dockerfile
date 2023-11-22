FROM golang:1.21-alpine as builder
ENV USER=appuser UID=10001
# Добавляем гит и софты для кафки gcc libc-dev librdkafka-dev pkgconf
RUN apk upgrade && apk add git gcc libc-dev librdkafka-dev pkgconf tzdata \
    # Добавляем юсера, чтобы им коннектится через второй стейдж scratch иначе у него не будет прав на запуск
    && adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR /go/src/pkg-outbox
COPY . .

RUN CGO_ENABLED=1 GOOS=linux \
    go build -o /compiled/migration \
    # Флаги нужны для добавления в бинарник cgo
    -a -ldflags '-linkmode external -extldflags "-static"' \
    # Флаг нужен для кафки
    -tags musl \
    cmd/migrate/main.go

# scratch легковесный образ, который в себе ничего не имеет и может только запустить бинарник
FROM busybox:musl
# переменные для считывания аргументов make команды
ARG ARGS
ARG C
ENV env_args=$ARGS
ENV env_c=$C
# копируем пароли для юсера и объявляем его
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER appuser:appuser
# копируем таймзону для запуска goose
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Копируем бинарник
COPY --from=builder /compiled/ /compiled/
# Копируем файлы миграций
COPY --from=builder /go/src/pkg-outbox/migration/migrations/ /compiled/migrations/

ENTRYPOINT /compiled/migration -dir=/compiled/migrations