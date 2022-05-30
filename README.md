## Install Prerequisites
Install Go, etcd and goreman:
```sh
sudo apt-get install golang
sudo apt-get install etcd
```

For goreman, download from https://github.com/mattn/goreman/releases/tag/v0.3.11 and put binaries under PATH.

## Initialize Standalone etcd Server and Service
```sh
etcd
```

## Initialize Multi-member Cluster
```sh
goreman -f Procfile start
```

This pulls up 3 etcd nodes. To stop one of them:
```sh
goreman run stop etcd1
```

## Run Server

Open another terminal and run the scripts:
```sh
go run server/server.go
```

## Run Client
```sh
go run client/client.go
```

## Regenerate gRPC Code
```sh
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    service/protoservices.proto
```

## Database APIs

You can do database management in ``client/client.go``.

The message struct ``RequestMsg`` has three fields: ``Operation``, ``Key``, ``Value``.

``Operation`` has 4 possible values: ``SetKV``, ``GetKey``, ``DeleteKey``, ``GetListValues``.

For ``SetKV``, it sets a key-value pair in the database. You need to provide ``Key``, ``Value`` fields in ``RequestMsg``.

For ``GetKey``, it returns the value of the requested key. You need to provide the field ``Key`` in ``RequestMsg``.

For ``DeleteKey``, it clears a key-value pair. You need to provide the field ``Key`` in ``RequestMsg``.

For ``GetListValues``, it returns a list of all values in the database.

Although ``Value`` passed to the server should be strings, you can define your own struct in ``client/client.go`` and stringify it before request.

See ``client/client.go`` for examples.