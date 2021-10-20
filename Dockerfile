FROM alpine

COPY kluster /usr/local/bin

ENTRYPOINT [ "customControllerCRD_CR" ]