FROM scratch

COPY nviz-linux .

ENTRYPOINT ["/nviz-linux"]
