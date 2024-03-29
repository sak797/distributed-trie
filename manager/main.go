package main

import (
	"../pb"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	// Argument parsing

	var managerPort int


	flag.IntVar(&managerPort, "managerPort", 3000,
		"Port on which manager should listen to client requests")

	flag.Parse()

	// Initialize the random number generator

	// Get hostname
	_, err := os.Hostname()
	if err != nil {
		// Without a host name we can't really get an ID, so die.
		log.Fatalf("Could not get hostname")
	}
	//TODO: Change later
	//name = "127.0.0.1"

	// Convert port to a string form
	portString := fmt.Sprintf(":%d", managerPort)
	// Create socket that listens on the supplied port
	c, err := net.Listen("tcp", portString)
	if err != nil {
		// Note the use of Fatalf which will exit the program after reporting the error.
		log.Fatalf("Could not create listening socket %v", err)
	}
	// Create a new GRPC server
	s := grpc.NewServer()

	// Initialize Manager

	manager := Manager{
		InitChan: make(chan PortIntroArgs),
		HeartbeatAckChan: make(chan HeartbeatAckArgs),
		SplitTrieRequestChan: make(chan SplitTrieRequestArgs),
		SplitWordsChan: make(chan SplitWordRequestArgs),
		SplitTrieRequestAckChan: make(chan SplitTrieRequestAckArgs),
	}


	pb.RegisterManagerServer(s, &manager)


	go manage(&manager)

	if err := s.Serve(c); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
	log.Printf("Done listening")
}