FROM golang:1.25.7 AS base 

WORKDIR /build

COPY go.mod go.sum ./ 

RUN go mod download 

COPY . .

RUN go build -o PloogleApiService

EXPOSE 9005

CMD ["/build/PloogleApiService"]
