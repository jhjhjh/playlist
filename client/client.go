package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func readInput(in *bufio.Reader) string {
	line, _ := in.ReadString('\n')
	return strings.ReplaceAll(line, "\n", "")
}

func statusPrint(messages <-chan string) {
	for message := range messages {
		fmt.Println(message)
		fmt.Println("Enter command [(p)lay, pa(u)se, (n)ext, p(r)evios, (a)dd, (d)elete, (e)xit, print(t), (s)ave, (l)oad]: ")
	}
}

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPlaylistServiceClient(conn)
	var str string
	in := bufio.NewReader(os.Stdin)
	messages := make(chan string)

	command := &pb.Command{Com: ""}
	_, err = c.ListFeatures(context.Background(), command)
	if err != nil {
		log.Fatalf("Stream error %v", err)
	}

	go statusPrint(messages)
	var name string
	var dur int
	//	var command *pb.Command
	for {
		command = nil
		messages <- "" //"Ener command [(p)lay, pa(u)se, (n)ext, p(r)evios, (a)dd, (d)elete, (e)xit]: "
		str = readInput(in)
		fmt.Println(str)
		switch str {
		case "exit", "e", "E":
			return
		case "add", "a", "A":
			fmt.Println("Enter song name: ")
			name = readInput(in)
			fmt.Println("Enter song duration in seconds")
			durStr := readInput(in)
			dur, _ = strconv.Atoi(durStr)
			command = &pb.Command{Com: "add", Name: name, Duration: int32(dur)}
		case "delete", "d", "D":
			fmt.Println("Enter song name")
			name = readInput(in)
			command = &pb.Command{Com: "delete", Name: name}
		case "play", "p", "P":
			command = &pb.Command{Com: "play"}
		case "pause", "u", "U":
			command = &pb.Command{Com: "pause"}
		case "next", "n", "N":
			command = &pb.Command{Com: "next"}
		case "previos", "r", "R":
			command = &pb.Command{Com: "previos"}
		case "save", "s", "S":
			command = &pb.Command{Com: "save"}
		case "load", "l", "L":
			command = &pb.Command{Com: "load"}
		case "t":
			command = &pb.Command{Com: "print"}

		}
		if command != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SendCommand(ctx, command)
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Response: %s", r)
		}
	}
}
