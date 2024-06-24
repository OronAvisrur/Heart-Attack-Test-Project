FROM alpine:latest

RUN mkdir /app

COPY knnApp /app

CMD [ "/app/knnApp" ]