FROM debian
RUN apt-get update
RUN apt-get -y upgrade
# Change TimeZone
ENV TZ=Europe/Moscow
# Clean APK cache
RUN rm -rf /var/cache/apk/*
WORKDIR /application
COPY main .
COPY ftp.json .
CMD ["./main" ]