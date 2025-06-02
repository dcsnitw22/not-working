FROM ubuntu:22.04

USER root

#ENV http_proxy="http://203.91.199.136:80/"
#ENV https_proxy="http://203.91.199.136:80/"
#ENV no_proxy="localhost,127.0.0.1,192.168.0.0/16,172.20.0.0/16,10.96.0.0/12,172.16.27.0/24,172.16.27.23,172.16.27.22,172.16.27.15,10.209.184.15,10.209.184.19,10.209.184.23"

#RUN apt-get update -y
#RUN apt-get install iputils-ping -y
#RUN apt-get install curl -y
#RUN apt-get install telnet -y
#RUN apt-get install dnsutils -y
#RUN mkdir -p /opt/w5gc/bin
# Copy binary
COPY ./bin/csp /bin

#temporary
#RUN mkdir -p /home/ubuntu/wipro5gc/testdata
RUN mkdir -p /root/wipro5gc/testdata
COPY testdata/smContextCreate.json testdata/smContextRelease.json testdata/smContextRetrieve.json testdata/smContextUpdate.json testdata/n1msgtest testdata/n2infoext1test testdata/n2infotest /root/wipro5gc/testdata
#COPY TestCases/UnitTestCases/pkg/smf/pdusmsp/sm/UpdateSmContextRequest/smContextUpdate_200.json /home/ubuntu/wipro5gc/testdata/smContextUpdate.json
#COPY testdata/smContextUpdate_swagger.json /home/ubuntu/wipro5gc/testdata/smContextUpdate.json
WORKDIR /opt/w5gc/bin
RUN export PATH=$PATH:/opt/w5gc/bin
# RUN export CONFIG_DIR=/root/go/src/w5gc.io/wipro5gcore/configs
CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
    exec csp
