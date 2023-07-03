FROM alpine

WORKDIR /

COPY ./app .

EXPOSE 8000

CMD ["./app"]