FROM alpine

RUN apk --update add ca-certificates

ADD dp_search_monitoring_linux_amd64 /

EXPOSE 8080

CMD /dp_search_monitoring_linux_amd64
