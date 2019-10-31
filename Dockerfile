FROM alpine:latest
WORKDIR /usr/src/app
EXPOSE 8081 8082 8083
COPY bin/simple-oauth-linux .
COPY client/templates/* ./client/templates/
COPY authorization/templates/* ./authorization/templates/
COPY protected/templates/* ./protected/templates/
ENTRYPOINT ["./simple-oauth-linux"]
CMD ["all"]