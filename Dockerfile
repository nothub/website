FROM busybox:stable-musl
COPY public /srv/http/hub.lol/
CMD ["httpd", "-f", "-u", "42000:42000", "-h", "/srv/http/hub.lol/"]
