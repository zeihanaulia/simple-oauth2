FROM alpine:latest
WORKDIR /usr/src/app
EXPOSE 8081 8082 8083
COPY bin/simple-oauth2 .
COPY client/templates/* ./client/templates/
COPY authorization/templates/* ./authorization/templates/
COPY protected/templates/* ./protected/templates/
ENTRYPOINT ["./simple-oauth2"]
CMD ["all"]