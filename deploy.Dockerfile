FROM alpine

WORKDIR /

COPY ./app .

EXPOSE 80
EXPOSE 443

CMD ["./app"]