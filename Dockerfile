FROM alpine:3

RUN apk add --no-cache     \
      "darkhttpd"          \
      "tini"               \
 && addgroup               \
      --gid "4242"         \
      --system             \
      "httpd"              \
 && adduser                \
      --uid "4242"         \
      --system             \
      --ingroup "httpd"    \
      --disabled-password  \
      --no-create-home     \
      "httpd"

COPY public /srv/http/hub.lol/

EXPOSE 8080

CMD ["/sbin/tini", "-vv", "--", "darkhttpd", "/srv/http/hub.lol/", "--chroot", "--uid",  "4242", "--gid",  "4242", "--port", "8080"]
