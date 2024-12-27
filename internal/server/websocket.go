package server

import (
	"fmt"
	"net/http"
	"sync"
	"math"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Server struct {
	Clients  map[string]*Client
	Mutex    sync.Mutex
	Upgrader websocket.Upgrader
}

func isWithinDistance(clientPos, objectPos [3]float64, maxDistance float64) bool {
    dx := clientPos[0] - objectPos[0]
    dy := clientPos[1] - objectPos[1]
    dz := clientPos[2] - objectPos[2]
    distance := math.Sqrt(dx*dx + dy*dy + dz*dz)
    return distance <= maxDistance
}

func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) Start(addr string) error {
    // Servir arquivos estáticos com suporte a CORS
    http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // Permite qualquer origem
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))).ServeHTTP(w, r)
    })

    // Endpoint WebSocket
    http.HandleFunc("/ws", s.HandleWebSocket)

    return http.ListenAndServe(addr, nil)
}

func filterObjects(world *World, clientPos [3]float64, maxDistance float64) map[string]Object {
    filtered := make(map[string]Object)
    for id, object := range world.Objects {
        if isWithinDistance(clientPos, object.Position, maxDistance) {
            filtered[id] = object
        }
    }
    return filtered
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := s.Upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Error upgrading connection:", err)
        return
    }
    defer conn.Close()

    client := &Client{
        ID:   r.RemoteAddr,
        Conn: conn,
    }

    s.Mutex.Lock()
    s.Clients[client.ID] = client
    s.Mutex.Unlock()

    defer func() {
        s.Mutex.Lock()
        delete(s.Clients, client.ID)
        s.Mutex.Unlock()
    }()

    // Posicione o cliente em uma posição inicial (ajuste conforme necessário)
    clientPos := [3]float64{0, 0, 0} // Exemplo de posição inicial

    world := NewWorld()
    filteredObjects := filterObjects(world, clientPos, 5.0) // Raio de 1 metro
    conn.WriteJSON(map[string]interface{}{
        "Objects": filteredObjects,
    })

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error reading message:", err)
            break
        }
        fmt.Println("Received:", string(msg))
    }
}