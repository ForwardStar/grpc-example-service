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
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	kv = clientv3.NewKV(cli)
)

// server is used to implement helloworld.GreeterServer.
type greeter_server struct {
	pb.UnimplementedGreeterServer
}

type etcd_server struct {
	pb.UnimplementedETCDServer
}

// SayHello implements helloworld.GreeterServer
func (s *greeter_server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *etcd_server) RequestETCD(ctx context.Context, in *pb.RequestMsg) (*pb.ResponseMsg, error) {
	log.Printf("Received request: %v", in.GetOperation())
	if in.GetOperation() == "SetKV" {
		_, err := kv.Put(ctx, in.GetKey(), in.GetValue())
		if err != nil {
			log.Fatalf("could not set: %v", err)
		}
		return &pb.ResponseMsg{Message: "Set successfully."}, nil
	}
	if in.GetOperation() == "GetKey" {
		r, err := kv.Get(ctx, in.GetKey())
		if err != nil {
			log.Fatalf("could not get: %v", err)
		}
		if len(r.Kvs) == 0 {
			return &pb.ResponseMsg{Message: "Could not find key " + in.GetKey()}, nil
		}
		return &pb.ResponseMsg{Message: string(r.Kvs[len(r.Kvs)-1].Value)}, nil
	}
	if in.GetOperation() == "DeleteKey" {
		_, err := kv.Delete(ctx, in.GetKey())
		if err != nil {
			log.Fatalf("could not delete: %v", err)
		}
		return &pb.ResponseMsg{Message: "Delete successfully."}, nil
	}
	if in.GetOperation() == "GetListValues" {

	}
	return &pb.ResponseMsg{Message: "Unknown type of operation " + in.GetOperation()}, nil
}

func main() {
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
	if err != nil {
		log.Fatalf("connect to etcd failed: %v", err)
	}
	log.Printf("server listening at %v", lis_etcd.Addr())
	s2 := grpc.NewServer()
	pb.RegisterETCDServer(s2, &etcd_server{})
	if err := s2.Serve(lis_etcd); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
