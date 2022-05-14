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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc-example-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr_greeting = flag.String("addr_greeting", "localhost:50051", "the address to connect to")
	addr_etcd     = flag.String("addr_etcd", "localhost:50052", "the address to connect to")
	name          = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	/*
		conn, err := grpc.Dial(*addr_greeting, grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		c := pb.NewGreeterClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	*/

	conn_etcd, err := grpc.Dial(*addr_etcd, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn_etcd.Close()

	c1 := pb.NewETCDClient(conn_etcd)

	ctx_etcd, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r1, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "SetKV", Key: "greeting", Value: "I am Haoxuan Xie"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r1.GetMessage())

	r2, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "GetKey", Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r2.GetMessage())

	r3, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "DeleteKey", Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r3.GetMessage())

	r4, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "GetKey", Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r4.GetMessage())

}
