FROM alpine

ADD dp_search_monitoring_linux_amd64 /

EXPOSE 8080

CMD /dp_search_monitoring_linux_amd64
