FROM golang:1.20 as base
#FROM ubuntu:22.04 as base
ENV http_proxy="http://10.250.110.36:8000/"
ENV https_proxy="http://10.250.110.36:8000/"
ENV no_proxy="localhost,127.0.0.1,192.168.0.0/16"

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
# Install csp
RUN make build
RUN make install
# Final image
#FROM scratch
#COPY --from=base /root/csp /root/go/bin/csp

WORKDIR /root/go/bin
#RUN export CONFIG_DIR=/root/go/src/w5gc.io/wipro5gcore/configs
CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
       exec csp

