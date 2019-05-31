FROM golang AS builder

WORKDIR /go/src/github.com/n0madic/datadog2image

ADD . .

RUN cd cmd/datadog2image/ && \
    go get -t && \
    go build


FROM chromedp/headless-shell

RUN ln -s /headless-shell/headless-shell /usr/bin/google-chrome

COPY --from=builder /go/src/github.com/n0madic/datadog2image/cmd/datadog2image/datadog2image /usr/bin/

ENTRYPOINT ["datadog2image"]