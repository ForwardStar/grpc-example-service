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
	"container/list"
	"context"
	"encoding/json"
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

// Replace DummyInfo with any struct
type DummyInfo struct {
	N int
	M int
}

func Stringify(data DummyInfo) string {

	stringified_data, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return string(stringified_data)

}

func Destringify(data string) any {

	if len(data) > 0 && data[0] == '[' {
		var stringified_data_list list.List
		in_json := false
		cut := ""
		for i := 1; i < len(data)-1; i++ {
			if data[i] == '{' {
				in_json = true
				cut += string(data[i])
				continue
			}
			if in_json {
				cut += string(data[i])
			}
			if in_json && data[i] == '}' {
				in_json = false
				stringified_data_list.PushBack(cut)
				cut = ""
			}
		}
		var destringified_data_list list.List
		for it := stringified_data_list.Front(); it != nil; it = it.Next() {
			value, ok := it.Value.(string)
			if ok {
				destringified_data_list.PushBack(Destringify(value))
			}
		}
		return destringified_data_list
	}

	var destringified_data DummyInfo

	err := json.Unmarshal([]byte(data), &destringified_data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return destringified_data

}

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

	r1, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "SetKV", Key: "greeting", Value: Stringify(DummyInfo{
		N: 6,
		M: 9,
	})})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r1.GetMessage())

	r2, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "GetKey", Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r2.GetMessage())

	r3, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "SetKV", Key: "greet", Value: Stringify(DummyInfo{
		N: 5,
		M: 10,
	})})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r3.GetMessage())

	r4, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "GetListValues"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r4.GetMessage())

	r5, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "DeleteKey", Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r5.GetMessage())

	r6, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "GetListValues"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r6.GetMessage())

	r7, err := c1.RequestETCD(ctx_etcd, &pb.RequestMsg{Operation: "DeleteKey", Key: "greet"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r7.GetMessage())

}
