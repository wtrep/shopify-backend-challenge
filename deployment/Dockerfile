FROM golang:1.15.1

ENV GOPRIVATE=github.com/wtrep/image-repo-backend
ARG GITHUB_TOKEN

WORKDIR /go/src/app

COPY backend .
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf "https://github.com/"
RUN go get -v -d
RUN go install -v

CMD ["image-repo-backend"]


