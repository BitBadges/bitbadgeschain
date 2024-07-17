FROM --platform=linux golang:1.21 AS builder

ENV COSMOS_VERSION=v0.47.5
RUN apt-get update && apt-get install -y git curl
RUN apt-get install -y make wget

WORKDIR /root

FROM golang:1.21 as v1

#Install dependencies
RUN apt-get update && apt-get install -y git curl
RUN apt-get install -y make wget


WORKDIR /home
RUN curl https://get.ignite.com/cli@v0.27.2 | bash
RUN mv ignite /usr/local/bin
RUN ignite version

#Clone the repository
# ENV VERSION=v1.0-betanet
# RUN git clone --branch ${VERSION} https://github.com/bitbadges/bitbadgeschain.git
RUN git clone https://github.com/bitbadges/bitbadgeschain.git

WORKDIR /home/bitbadgeschain

RUN ignite chain build --skip-proto

WORKDIR /

ENV LOCAL=/usr/local
ENV DAEMON_NAME=bitbadgeschaind
ENV DAEMON_HOME=/root/.bitbadgeschain

RUN mv /go/bin/bitbadgeschaind  ${LOCAL}/bin/bitbadgeschaind

EXPOSE 26656 26657 26660 6060 9090 1317

ENTRYPOINT [ "bitbadgeschaind" ]