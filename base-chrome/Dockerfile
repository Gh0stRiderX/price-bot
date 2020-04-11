FROM arm32v7/debian:latest

RUN \
    apt-get update -y \
    && apt-get install chromium-browser -y \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
    && ln -s /usr/bin/chromium-browser /usr/bin/google-chrome


ENV PATH /usr/bin/google-chrome:$PATH
