package server

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (s *Server) Broadcast(message Message) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	for _, client := range s.Clients {
		err := client.Conn.WriteMessage(1, data)
		if err != nil {
			fmt.Println("Error broadcasting message:", err)
		}
	}
}
