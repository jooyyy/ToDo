FROM alpine

COPY ./bin/todo .

COPY config.json .

ENTRYPOINT ["./todo"]

