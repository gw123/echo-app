FROM alpine:3.12
ENV TZ Asia/Shanghai

RUN apk add --no-cache \
        tzdata \
        && rm -f /etc/localtime \
        && ln -s /usr/share/zoneinfo/$TZ /etc/localtime \
        && echo $TZ > /etc/timezone

RUN apk add --update --no-cache \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY resources/views /usr/local/var/echoapp/resources/views
COPY resources/public /usr/local/var/echoapp/resources/public
COPY echoapp /usr/local/bin/echoapp
#COPY docker-entrypoint /usr/local/bin/

WORKDIR /usr/local/var/echoapp
EXPOSE 80
#VOLUME ["/opt", "/usr/local/etc/echoapp"]
#ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]