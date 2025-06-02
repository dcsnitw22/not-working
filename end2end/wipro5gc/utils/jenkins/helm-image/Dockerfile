FROM jenkins/jenkins:latest
#SHELL ["/bin/bash", "-c"]

USER root

ENV http_proxy=http://10.250.110.36:8000/
ENV https_proxy=http://10.250.110.36:8000/
ENV no_proxy=localhost,127.0.0.1,192.168.0.0/16,10.250.110.39,10.96.0.0/12

# Install all the dev dependencies
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y wget
RUN apt-get install -y nano
RUN apt-get install -y unzip