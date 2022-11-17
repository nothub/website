FROM nixos/nix

RUN nix-channel --update  \
 && nix-env --upgrade \
 && nix-env --install --attr nixpkgs.darkhttpd

COPY public /srv/http/hub.lol/

ENTRYPOINT ["darkhttpd", "/srv/http/hub.lol/", "--chroot", "--no-listing"]
CMD ["--port", "8080"]
