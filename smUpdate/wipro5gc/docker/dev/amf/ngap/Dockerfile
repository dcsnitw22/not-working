FROM golang:1.20 as base
#FROM ubuntu:22.04 as base
ENV http_proxy="http://203.91.199.136:80/"
ENV https_proxy="http://203.91.199.136:80/"
ENV no_proxy="localhost,127.0.0.1,172.16.27.10,172.16.27.9,172.16.27.12,172.16.27.21,172.16.27.11,172.20.135.33,172.16.27.14,10.98.249.67,10.100.161.68,10.110.163.12,172.16.27.3,172.16.27.1,172.16.27.2,172.16.27.4,172.16.27.7,gitlab.wipro5gc.com,10.209.184.15,10.209.184.50,10.102.118.32,10.97.178.232,10.99.59.146,172.16.27.8,lb-apiserver.kubernetes.local,192.168.0.0/16,10.96.0.0/12,10.209.184.19,10.209.184.23,10.100.111.77,10.96.0.0/12,172.20.0.0/16,10.209.184.109,10.209.184.105,10.209.184.113,172.16.27.21,172.16.27.22,172.16.27.23,172.16.27.25,10.209.184.26,10.209.184.38,10.108.129.34"

#RUN mkdir /opt/w5gc
RUN mkdir -p /root/go/src/w5gc.io/wipro5gcore
RUN apt-get update && apt-get install -y --no-install-recommends \
                # general tools
                inetutils-traceroute \
                iproute2 \
                iputils-ping \
		make \
        && rm -rf /var/lib/apt/lists/*
# Copy code
COPY . /root/go/src/w5gc.io/wipro5gcore
WORKDIR /root/go/src/w5gc.io/wipro5gcore
# Install ngap
RUN make build
RUN make install
# Final image
#FROM scratch
#COPY --from=base /root/ngap /root/go/bin/ngap

WORKDIR /root/go/bin
#RUN export CONFIG_DIR=/root/go/src/w5gc.io/wipro5gcore/configs
CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
       exec ngap

