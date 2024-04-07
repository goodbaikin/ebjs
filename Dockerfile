FROM ubuntu:jammy
RUN apt-get update &&\
  apt-get install -y \
    git \
    autoconf \
    libtool \
    make \
    meson \
    gcc \
    nasm \
    cmake \
    g++ \
    pkg-config \
    libdevil-dev \
    intel-media-va-driver-non-free \
    libva-dev \
    libmfx-dev \
    libx264-dev \
    libfdk-aac-dev \
    libvpx-dev \
    libvorbis-dev

RUN git clone https://github.com/AviSynth/AviSynthPlus.git &&\
  cd AviSynthPlus &&\
  mkdir build &&\
  cd build &&\
  cmake .. &&\
  make &&\
  make install PREFIX=/usr/local

RUN git clone --depth 1 https://git.ffmpeg.org/ffmpeg.git &&\
  cd ffmpeg &&\
  ./configure \
    --enable-gpl \
    --enable-avisynth \
    --enable-libmfx \
    --enable-libx264 \
    --enable-libvpx \
    --enable-libvorbis \
    --enable-libfdk-aac \
    --enable-nonfree \
    &&\
  make &&\
  make install

RUN git clone https://github.com/l-smash/l-smash.git &&\
  cd l-smash &&\
  ./configure --enable-shared &&\
  make &&\
  make install

RUN git clone https://github.com/HomeOfAviSynthPlusEvolution/L-SMASH-Works.git &&\
  cd L-SMASH-Works/AviSynth &&\
  LDFLAGS="-Wl,-Bsymbolic" meson build &&\
  cd build &&\
  ninja &&\
  ninja install &&\
  cd ../../.. &&\
  ldconfig

RUN git clone https://github.com/tobitti0/chapter_exe.git &&\
  cd chapter_exe/src &&\
  make &&\
  cp chapter_exe /usr/local/bin/
RUN git clone https://github.com/tobitti0/logoframe.git &&\
  cd logoframe/src &&\
  make &&\
  cp logoframe /usr/local/bin/
RUN git clone https://github.com/tobitti0/join_logo_scp.git &&\
  cd join_logo_scp/src &&\
  make &&\
  cp join_logo_scp /usr/local/bin/
RUN git clone https://github.com/tobitti0/delogo-AviSynthPlus-Linux.git &&\
  cd delogo-AviSynthPlus-Linux/src &&\
  sed -i -e 's/^LDLAGS = -shared -fPIC/LDLAGS = -shared -fPIC -lstdc++/' Makefile &&\
  make &&\
  make install

FROM ubuntu:jammy
COPY --from=0 /usr/local /usr/local
RUN ldconfig
RUN apt-get update &&\
  apt-get install --no-install-recommends -y \
    libdevil1c2 \
    libvpx7 \
    libx264-163 \
    libfdk-aac2 \
    intel-media-va-driver-non-free \ 
    libmfx1 \
    libva2 \
    libva-drm2 \
    libva-x11-2 \
    &&\
  apt-get clean &&\
  rm -rf /var/lib/apt/lists
ADD bin/ebjs /ebjs
ENTRYPOINT ["/ebjs"]