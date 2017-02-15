FROM alpine:latest

COPY dumbledore /usr/local/bin
COPY *.so /usr/lib/dumbledore/plugins/

EXPOSE 8080

VOLUME [ "/etc/dumbledore" ]

ENTRYPOINT [ "dumbledore" ]
CMD [ "-tcp", ":8080", "-f", "/etc/dumbledore/example.conf", "-plugin-dir", "/usr/lib/dumbledore/plugins/" ]
