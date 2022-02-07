FROM alpine

COPY ./bin/todo .
COPY config.json .

# set timezone to +8.0
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

ENTRYPOINT ["./todo"]

