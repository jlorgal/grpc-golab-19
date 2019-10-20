Compile proto file according to guidelines for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). It will generate the file `compose/proto/compose.pb.go`:

```sh
protoc -I. \
       -I/usr/local/include \
       -I$GOPATH/src \
       -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --go_out=plugins=grpc:. \
       *.proto
```

Generate reverse proxy (`compose/proto/compose.pb.gw.go`):

```
protoc -I. \
       -I/usr/local/include \
       -I$GOPATH/src \
       -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --grpc-gateway_out=logtostderr=true:. \
       *.proto
```
