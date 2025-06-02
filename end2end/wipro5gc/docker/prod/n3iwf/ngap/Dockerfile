FROM custom_ubuntu:latest

USER root

COPY ./n3iwf-ngap /bin

WORKDIR /opt/w5gc/bin
RUN export PATH=$PATH:/opt/w5gc/bin
CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
    exec n3iwf-ngap