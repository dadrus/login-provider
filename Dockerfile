FROM golang:buster
LABEL maintainer="dimitrij.drus@innoq.com"

RUN useradd --system --user-group login-provider

ADD login-provider /opt/login-provider/login-provider
ADD web /opt/login-provider/web

RUN chown -R login-provider:login-provider /opt/login-provider

USER login-provider

ENV GIN_MODE release

WORKDIR /opt/login-provider
ENTRYPOINT ["./login-provider"]

EXPOSE 8080