FROM alpine:3.18

ARG BINARY

ENV STRATE_GO_BIN=/opt/strate-go/bin/strate-go

COPY $BINARY $STRATE_GO_BIN

RUN apk add --no-cache libc6-compat gcompat && \
    chmod a+x $STRATE_GO_BIN && \
    ln -s $STRATE_GO_BIN /bin/strate-go

ENTRYPOINT ["strate-go"]
