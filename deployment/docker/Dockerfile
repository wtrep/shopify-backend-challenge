FROM golang:1.15.1

ENV GOPRIVATE=github.com/wtrep/image-repo-backend

WORKDIR /go/src/app

COPY backend ./backend
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go get -v -d
RUN go install -v

CMD ["image-repo-backend"]


