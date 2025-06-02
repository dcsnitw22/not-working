FROM ubuntu:22.04

USER root

# ENV http_proxy="http://203.91.199.136:80/"
# ENV https_proxy="http://203.91.199.136:80/"
# ENV no_proxy="localhost,127.0.0.1,192.168.0.0/16,172.16.27.15,10.209.184.15,10.209.184.19,10.209.184.23,.svc.cluster.local"

RUN mkdir -p /opt/w5gc/bin
#RUN mkdir -p /home/ubuntu/wipro5gc/n1n2data
RUN mkdir -p /root/wipro5gc/n1n2data
# Copy binary
COPY ./bin/pdusmsp /bin
#temporary
# RUN mkdir -p /home/requests
# COPY smContextCreate.json smContextRelease.json smContextUpdate.json n1msgtest n2infoext1test n2infotest /home/requests/

COPY ./bin/upfgw /bin


WORKDIR /opt/w5gc/bin
RUN export PATH=$PATH:/opt/w5gc/bin
# RUN export CONFIG_DIR=/root/go/src/w5gc.io/wipro5gcore/configs
#CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
#    exec pdusmsp && \
#    exec upfgw  
    # --request-type="update"

CMD ["/bin/sh", "-c", "rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && exec pdusmsp & exec upfgw "]
