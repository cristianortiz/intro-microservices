
# build a tiny docker image, to copy over the executable
FROM alpine:latest

RUN mkdir /app

COPY mailServiceApp /app
COPY templates /templates

CMD [ "/app/mailServiceApp" ]