FROM openjdk:8-jdk-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk add dumb-init libc6-compat dbus curl \
    && apk add --update --no-cache ttf-dejavu fontconfig \
    && rm -rf /var/cache/apk/*

ARG app=ygoapi
ARG path=/home/${app}
RUN mkdir -p ${path}
COPY files ${path}/files
COPY ygoapi ${path}
COPY application.yml ${path}
WORKDIR ${path}

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./ygoapi"]
