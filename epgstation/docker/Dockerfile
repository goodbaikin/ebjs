FROM l3tnun/epgstation:master-debian
RUN apt-get update &&\
  apt-get install -y --no-install-recommends ffmpeg &&\
  apt-get clean &&\
  rm -rf /var/lib/apt/lists

ADD epgstation/docker/amatsukaze-add-task /usr/local/bin
ADD epgstation/docker/amatsukaze-add-task.sh /usr/local/bin
ADD bin/ebjclient /usr/local/bin