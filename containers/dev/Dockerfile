# workspace (GOPATH) configured at /go
FROM golang:1.21 as builder

RUN go install github.com/cespare/reflex@v0.3.1

COPY reflex.conf /
COPY .netrc /root/.netrc

RUN go env -w GOPRIVATE=gitlab.iman.uz,gitlab.iman.uz/*
RUN export GOPRIVATE=gitlab.iman.uz,gitlab.iman.uz/*
RUN export GONOSUMDB=gitlab.iman.uz,gitlab.iman.uz/*
RUN export GONOPROXY=gitlab.iman.uz,gitlab.iman.uz/*

COPY start.sh /

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]
