FROM golang:latest
#SHELL ["/bin/bash", "-c"]

ENV http_proxy=http://10.250.110.36:8000/
ENV https_proxy=http://10.250.110.36:8000/
ENV no_proxy=localhost,127.0.0.1,192.168.0.0/16,10.250.110.39,10.96.0.0/12

RUN apt update
RUN apt install -y ca-certificates curl gnupg

RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
RUN chmod a+r /etc/apt/keyrings/docker.gpg

RUN echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt update
RUN apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Install all the dev dependencies
RUN apt install -y build-essential
RUN apt install -y wget
RUN apt install -y nano

#RUN wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
#RUN echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
RUN export PATH=$PATH:/usr/local/go/bin

RUN apt update
RUN apt install -y default-jdk
RUN which java

#go env -w GO111MODULE=off
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest

RUN apt update
RUN apt install -y python3.11
RUN apt install -y python3-pip
RUN pip3 install robotframework --break-system-packages
RUN pip3 install requests --break-system-packages
RUN pip3 install robotframework-requests --break-system-packages
RUN pip3 install robotframework-jsonlibrary --break-system-packages
RUN pip3 install robotframework-lint --break-system-packages
RUN pip3 install robotframework-seleniumlibrary --break-system-packages
RUN pip3 install pyyaml --break-system-packages