FROM golang:1.21


RUN echo "deb http://packages.cloud.google.com/apt gcsfuse-buster main" > /etc/apt/sources.list.d/gcsfuse.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

RUN apt-get update && \
    apt-get install -y gcsfuse 

RUN mkdir /mnt/gcs-bucket-test1312
RUN gcsfuse --implicit-dirs bucket-test1312 /mnt/gcs-bucket-test1312

WORKDIR /go/src/app
COPY . .

CMD ["go", "run", "main.go"]


