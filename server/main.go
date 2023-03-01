package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	pb "server/pb"
	"time"

	"google.golang.org/grpc"
)

type State int

const (
	Stop  State = 0
	Play  State = 1
	Pause State = 2
)

type Command struct {
	com      string
	name     string
	duration int
}

var head *PlayListElement
var current *PlayListElement
var state State

var commands chan *pb.Command
var messages chan string

type PlayListElement struct {
	name     string
	duration int
	next     *PlayListElement
	prev     *PlayListElement
}

func add(name string, dur int) {

	if head == nil {
		head = &PlayListElement{name: name, duration: dur}
		current = head
	} else {
		end := head
		for end.next != nil {
			end = end.next
		}
		end.next = &PlayListElement{name: name, duration: dur}
		end.next.prev = end
		end = end.next
	}
}

func del(name string) error {
	if head == nil {
		return errors.New("Playlsit is empty")
	}
	tmp := head
	for tmp != nil {
		if tmp.name == name {
			if tmp == current {
				if state == Play || state == Pause {
					return errors.New("Cannot delete current song ")
				} else {
					current = tmp.next
				}
			}
			if tmp == head {
				head = head.next
				tmp.next = nil
				tmp = nil
				if head != nil {
					head.prev = nil
				}
			} else {
				if tmp.next == nil {
					tmp.prev.next = nil
					tmp.prev = nil
					tmp = nil
				} else {
					tmpprev := tmp.prev
					tmpnext := tmp.next
					tmp.next = nil
					tmp.prev = nil
					tmpprev.next = tmpnext
					tmpnext.prev = tmpprev
					tmp = nil
				}
			}
			if current == nil {
				current = head
			}
		} else {
			tmp = tmp.next
		}
	}
	return nil
}

func print() error {
	if head == nil {
		return errors.New("Playlist is empty")
	}
	tmp := head
	for tmp != nil {
		fmt.Print("Song name:\t" + tmp.name + "\tduration:\t")
		fmt.Print(tmp.duration)
		if tmp.prev != nil {
			fmt.Print("\tPrevios song:\t" + tmp.prev.name)
		}
		if tmp.next != nil {
			fmt.Print("\tNext song:\t" + tmp.next.name)
		}
		if tmp == current {
			fmt.Print("*")
		}
		fmt.Println()
		tmp = tmp.next
	}
	return nil
}

func statusPrint(messages <-chan string) {
	for message := range messages {
		fmt.Println(message)
	}
}

func control(messages chan<- string, commands chan *pb.Command) {
	playCommands := make(chan Command)
	go play(messages, commands, playCommands)
	ch := make(chan bool, 1)
	for {
		select {
		case command := <-commands:
			switch command.GetCom() {
			case "add":
				ch <- true
				add(command.GetName(), int(command.GetDuration()))
				<-ch
			case "delete":
				ch <- true
				err := del(command.GetName())
				if err != nil {
					messages <- err.Error()
				}
				<-ch
			case "play":
				ch <- true
				state = Play
				cmd := Command{com: "play"}
				playCommands <- cmd
				<-ch
			case "pause":
				ch <- true
				switch state {
				case Play:
					state = Pause
					messages <- "Pause"
				case Pause:
					state = Play
					messages <- "Continue"
				}
				<-ch
			case "next":
				ch <- true
				if head == nil {
					messages <- "Empty Playlist"
				} else if current.next != nil {
					current = current.next
					cmd := Command{com: "play"}
					playCommands <- cmd
					state = Play
				} else {
					current = head
					messages <- "Playlist end"
				}
				<-ch
			case "previos":
				ch <- true
				if head == nil {
					messages <- "Empty Playlist"
				} else if current.prev != nil {
					current = current.prev
					cmd := Command{com: "play"}
					playCommands <- cmd
					state = Play
				}
				<-ch
			case "print":
				print()
			}
		}
	}
}

func play(messages chan<- string, commands chan<- *pb.Command, changeSong <-chan Command) {
	ticker := time.NewTicker(time.Second)
	var duration int
	for {
		select {
		case <-ticker.C:
			switch state {
			case Play:
				duration--
				if duration%10 == 0 {
					messages <- fmt.Sprintf("Song:\t%s\tlast:\t%d", current.name, duration)

				}
				if duration == 0 {
					command := &pb.Command{Com: "next"}
					state = Stop
					commands <- command
				}
			}
		case command := <-changeSong:
			switch command.com {
			case "play":
				if current != nil {
					duration = current.duration
					messages <- fmt.Sprintf("Song:\t%s\tduration:\t%d", current.name, current.duration)
					state = Play
				} else {
					state = Stop
					messages <- "Cannot play empty list"
				}
			}
		}
	}

}

type PlaylistServiceServer struct {
	pb.UnimplementedPlaylistServiceServer
}

func (server *PlaylistServiceServer) SendCommand(ctx context.Context, in *pb.Command) (*pb.Response, error) {
	log.Printf("Recieved: %v", in.GetCom())
	commands <- in
	//messages <- in.GetCom()
	return &pb.Response{Data: "Command " + in.GetCom() + " was resieved by server"}, nil
}

func main() {
	messages = make(chan string)
	commands = make(chan *pb.Command)
	go statusPrint(messages)
	go control(messages, commands)
	state = Stop

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPlaylistServiceServer(s, &PlaylistServiceServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
