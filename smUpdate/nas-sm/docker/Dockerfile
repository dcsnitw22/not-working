FROM ubuntu:22.04

USER root

#ENV http_proxy="http://203.91.199.136:80/"
#ENV https_proxy="http://203.91.199.136:80/"
#ENV no_proxy="localhost,127.0.0.1,192.168.0.0/16,172.16.27.15,172.16.27.22,172.16.27.23,10.209.184.15,10.209.184.19,10.209.184.23,172.20.0.0/16,10.96.0.0/12,172.16.27.0/24,.svc.cluster.local,.SVC"

#RUN apt-get update -y
#RUN apt-get install iputils-ping -y
#RUN apt-get install curl -y
#RUN apt-get install telnet -y
#RUN mkdir -p /opt/w5gc/bin
# Copy binary
COPY ./bin/grpcNasSmf /bin

WORKDIR /opt/w5gc/bin

RUN export PATH=$PATH:/opt/w5gc/bin

EXPOSE 50052

CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
   exec grpcNasSmf
