FROM ubuntu:latest
# default answers to all the questions
ENV DEBIAN_FRONTEND noninteractive 

RUN apt-get update && \
    apt-get -y --no-install-recommends install \
    build-essential \
    curl \
    nginx \
    libffi-dev \
    golang \
    git \
    python3 \
    python3-dev \
    python3-setuptools \
    python3-pip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*


WORKDIR /
ENV GOPATH /go
ENV PATH ${PATH}:/steres

COPY requirements.txt steres/requirements.txt
RUN pip3 install --no-cache-dir -r steres/requirements.txt

COPY steres volume steres/
COPY src/*.go steres/src/
COPY tools/* steres/tools/
# steres is cool
WORKDIR /steres


