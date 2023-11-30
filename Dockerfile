FROM gh0st42/coreemu7:1.0.0
LABEL Description="CORE Docker Ubuntu Image"

WORKDIR /home/torrent

# install go

COPY --from=golang:1.21.1 /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

# copy project files

COPY ./bin ./bin
# build project


# RUN	go build -o ./bin/node ./src/node
# RUN	go build -o ./bin/tracker ./src/tracker

# COPY ./topologies/. /usr/local/lib/core/.

RUN echo 'core-daemon > /var/log/core-daemon.log 2>&1 & \n\
sleep 1 \n\
core-gui' > /root/entryPoint.sh

ENTRYPOINT ["/bin/sh", "-c", "\"/root/entryPoint.sh\""]