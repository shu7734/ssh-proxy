FROM golang:1.7-alpine

RUN set -ex && \
  apk --update --upgrade add \
  git && \
  go get github.com/Masterminds/glide && \
  find / -type f -iname \*.apk-new -delete && \
  rm -rf /var/cache/apk/*

CMD ["ash"]
