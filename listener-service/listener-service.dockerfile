# build a tiny docker image, to copy over the executable
FROM alpine:latest

RUN mkdir /app

COPY listenerApp /app

CMD [ "/app/listenerApp" ]