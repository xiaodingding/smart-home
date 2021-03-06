// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package stream

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

type Hub struct {
	sessions    map[*Client]bool
	subscribers map[string]func(client IStreamClient, msg Message)
	sync.Mutex
	broadcast chan []byte
	interrupt chan os.Signal
}

func NewHub() *Hub {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	hub := &Hub{
		sessions:    make(map[*Client]bool),
		broadcast:   make(chan []byte, maxMessageSize),
		subscribers: make(map[string]func(client IStreamClient, msg Message)),
		interrupt:   interrupt,
	}
	go hub.Run()

	return hub
}

func (h *Hub) AddClient(client *Client) {

	defer func() {
		h.Lock()
		delete(h.sessions, client)
		h.Unlock()
		log.Infof("websocket session from ip: %s closed", client.Ip)
	}()

	h.Lock()
	h.sessions[client] = true
	h.Unlock()

	log.Infof("new websocket xsession established, from ip: %s", client.Ip)

	for {
		op, r, err := client.Connect.NextReader()
		if err != nil {
			break
		}
		switch op {
		case websocket.TextMessage:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				break
			}
			h.Recv(client, message)
		}
	}
}

func (h *Hub) Run() {

	for {
		select {
		case m := <-h.broadcast:
			h.Lock()
			for client := range h.sessions {
				client.Send <- m
			}
			h.Unlock()
		case <-h.interrupt:
			//fmt.Println("Close websocket client session")
			h.Lock()
			for client := range h.sessions {
				client.Close()
				delete(h.sessions, client)
			}
			h.Unlock()
		}

	}
}

func (h *Hub) Recv(client *Client, b []byte) {

	//fmt.Printf("client(%v), message(%v)\n", client, string(b))

	msg, err := NewMessage(b)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch msg.Command {
	case "client_info":
		client.UpdateInfo(msg)

	default:

		h.Lock()
		f, ok := h.subscribers[msg.Command]
		h.Unlock()
		if ok {
			f(client, msg)
		}
	}
}

func (h *Hub) Send(client *Client, message []byte) {
	client.Send <- message
}

func (h *Hub) Broadcast(message []byte) {
	h.Lock()
	h.broadcast <- message
	h.Unlock()
}

func (h *Hub) Clients() (clients []*Client) {

	clients = []*Client{}
	for c := range h.sessions {
		clients = append(clients, c)
	}

	return
}

func (h *Hub) Subscribe(command string, f func(client IStreamClient, msg Message)) {
	log.Infof("subscribe %s", command)
	h.Lock()
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.subscribers[command] = f
	h.Unlock()
}

func (h *Hub) UnSubscribe(command string) {
	h.Lock()
	if h.subscribers[command] != nil {
		delete(h.subscribers, command)
	}
	h.Unlock()
}

func (h *Hub) Subscriber(command string) (f func(client IStreamClient, msg Message)) {
	h.Lock()
	if h.subscribers[command] != nil {
		f = h.subscribers[command]
	}
	h.Unlock()
	return
}
