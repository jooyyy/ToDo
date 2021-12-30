FROM alpine

COPY ./bin/todo .
COPY config.json .

# set timezone to +8.0
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

ENTRYPOINT ["./todo"]

