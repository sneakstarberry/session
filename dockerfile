FROM golang:latest


WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .



RUN go mod download
RUN go get github.com/pilu/fresh

EXPOSE 8888

COPY . /app

CMD ["fresh"]