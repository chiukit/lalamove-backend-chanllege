FROM golang:1.8

RUN mkdir -p /go/src/lalamove
WORKDIR /go/src/lalamove

ENV GOOGLE_API_KEY=[Your Google API Key]

COPY . /go/src/lalamove

RUN go get
RUN go build

ENTRYPOINT ["/go/src/lalamove/lalamove"]

EXPOSE 3000