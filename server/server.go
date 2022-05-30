/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-example-service/service"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379", "http://127.0.0.1:2380"},
		DialTimeout: 5 * time.Second,
	})
	kv = clientv3.NewKV(cli)
)

// server is used to implement helloworld.GreeterServer.
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

type etcdServer struct {
	pb.UnimplementedETCDWrapperServer
}

// SayHello implements helloworld.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *etcdServer) SetKV(ctx context.Context, in *pb.SetKVRequest) (*pb.SetKVResponse, error) {
	log.Printf("Received request: SetKV, Key: " + in.GetKey() + ", Value: " + in.GetValue().String())
	jsonBytes, err := json.Marshal(in.GetValue())
	if err != nil {
		log.Fatalf("Could not encode: %v", err)
	}
	_, err = kv.Put(ctx, "/mytest/"+in.GetKey(), string(jsonBytes))
	if err != nil {
		log.Fatalf("Could not set: %v", err)
	}
	return &pb.SetKVResponse{Message: "OK"}, nil
}

func (s *etcdServer) GetKey(ctx context.Context, in *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	log.Printf("Received request: GetKey, Key: " + in.GetKey())
	rawData, err := kv.Get(ctx, "/mytest/"+in.GetKey())
	if err != nil {
		log.Fatalf("Could not get: %v", err)
	}
	if len(rawData.Kvs) > 0 {
		var value pb.DummyInfo
		errorMsg := json.Unmarshal(rawData.Kvs[len(rawData.Kvs)-1].Value, &value)
		if errorMsg != nil {
			log.Printf("Could not decode: %v", errorMsg)
		}
		return &pb.GetKeyResponse{Value: &value}, nil
	} else {
		return &pb.GetKeyResponse{Value: nil}, nil
	}
}

func (s *etcdServer) DeleteKey(ctx context.Context, in *pb.DeleteKeyRequest) (*pb.DeleteKeyResponse, error) {
	log.Printf("Received request: DeleteKey")
	_, err := kv.Delete(ctx, "/mytest/"+in.GetKey())
	if err != nil {
		log.Fatalf("Could not delete: %v", err)
	}
	return &pb.DeleteKeyResponse{Message: "OK"}, nil
}

func (s *etcdServer) ListValues(ctx context.Context, in *pb.ListValuesRequest) (*pb.ListValuesResponse, error) {
	log.Printf("Received request: GetListValues")
	rawData, err := kv.Get(ctx, "/mytest/", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Could not get: %v", err)
	}
	if len(rawData.Kvs) > 0 {
		values := make([]pb.DummyInfo, len(rawData.Kvs))
		valuePointers := make([]*pb.DummyInfo, len(rawData.Kvs))
		for i := range values {
			valuePointers[i] = &values[i]
		}
		for i := range rawData.Kvs {
			errorMsg := json.Unmarshal(rawData.Kvs[i].Value, valuePointers[i])
			if errorMsg != nil {
				log.Printf("Could not decode: %v", errorMsg)
			}
		}
		return &pb.ListValuesResponse{Values: valuePointers}, nil
	} else {
		return &pb.ListValuesResponse{Values: nil}, nil
	}
}

func main() {
	if err != nil {
		log.Fatalf("connect to etcd failed: %v", err)
	}
	flag.Parse()
	/*
		lis_greeting, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("server listening at %v", lis_greeting.Addr())
		s1 := grpc.NewServer()
		pb.RegisterGreeterServer(s1, &greeter_server{})
		if err := s1.Serve(lis_greeting); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	*/

	lis_etcd, err := net.Listen("tcp", fmt.Sprintf(":%d", *port+1))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis_etcd.Addr())
	s2 := grpc.NewServer()
	pb.RegisterETCDWrapperServer(s2, &etcdServer{})
	if err := s2.Serve(lis_etcd); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
