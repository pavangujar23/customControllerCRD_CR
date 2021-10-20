FROM alpine

COPY customControllerCRD_CR /usr/local/bin

ENTRYPOINT [ "customControllerCRD_CR" ]