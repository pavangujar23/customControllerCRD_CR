FROM alphine

COPY kluster /usr/local/bin

ENTRYPOINT [ "kluster" ]