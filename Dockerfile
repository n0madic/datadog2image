FROM golang AS builder

WORKDIR /go/src/github.com/n0madic/datadog2image

ADD . .

RUN cd cmd/datadog2image/ && \
    go get -t && \
    go install -ldflags="-s -w"


FROM chromedp/headless-shell

RUN ln -s /headless-shell/headless-shell /usr/bin/google-chrome

RUN apt-get update -qq && \
    apt-get remove tzdata -yqq && \
    apt-get install dumb-init tzdata -yqq && \
	apt-get autoremove -yqq && \
	rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/* /usr/bin/

ENTRYPOINT ["dumb-init", "--", "datadog2image"]
