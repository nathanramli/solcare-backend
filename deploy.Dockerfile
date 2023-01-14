FROM alpine

WORKDIR /

COPY ./app .

EXPOSE 80

CMD ["./app"]