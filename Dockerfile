FROM golang:1.13.4 AS build
WORKDIR /src
RUN go get -u github.com/golang/glog
COPY . /src
RUN CGO_ENABLED=0 go build -a main.go

FROM scratch AS final
WORKDIR /app
COPY --from=build /src/main /app/
ENTRYPOINT ["/app/main"]
