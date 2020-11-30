FROM alpine:latest 
RUN mkdir -p /server/static
COPY webserver /server
COPY static /server/static
RUN apk --update add git less openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/* && \
    mkdir /munki_repo
WORKDIR /server
ENTRYPOINT ["/server/webserver"] 
EXPOSE 8080