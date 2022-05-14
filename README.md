## Install Prerequisites
Install Go and etcd.

## Initialize etcd Server and Service
```sh
etcd
```

Open another terminal and run the scripts:
```sh
go run server/server.go
```

## Run Client
```sh
go run client/client.go
```

## Database APIs

You can do database management in ``client/client.go``.

The message struct ``RequestMsg`` has three fields: ``Operation``, ``Key``, ``Value``.

``Operation`` has 4 possible values: ``SetKV``, ``GetKey``, ``DeleteKey``, ``GetListValues``.

For ``SetKV``, it sets a key-value pair in the database. You need to provide ``Key``, ``Value`` fields in ``RequestMsg``.

For ``GetKey``, it returns the value of the requested key. You need to provide the field ``Key`` in ``RequestMsg``.

For ``DeleteKey``, it clears a key-value pair. You need to provide the field ``Key`` in ``RequestMsg``.

For ``GetListValues``, it returns a list of all values in the database.

See ``client/client.go`` for examples.