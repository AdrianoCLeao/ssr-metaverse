package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func filterObjects(world *World, clientPos [3]float64, maxDistance float64) map[string]Object {
	filtered := make(map[string]Object)
	for id, object := range world.Objects {
		if isWithinDistance(clientPos, object.Position, maxDistance) {
			filtered[id] = object
		}
	}
	return filtered
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
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.StaticFS("/assets", http.Dir("./assets"))

	router.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from RESTful API!"})
	})

	router.GET("/ws", func(c *gin.Context) {
		s.HandleWebSocket(c.Writer, c.Request)
	})

	return router.Run(addr)
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

	world := NewWorld()

	prevFilteredObjects := make(map[string]Object)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		var payload struct {
			Position [3]float64 `json:"position"`
			Rotation [3]float64 `json:"rotation"`
		}
		if err := json.Unmarshal(msg, &payload); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		maxDistance := 50.0
		filteredObjects := filterObjects(world, payload.Position, maxDistance)

		changes := map[string]interface{}{
			"added":   make(map[string]Object),
			"removed": make([]string, 0),
		}

		for id, object := range filteredObjects {
			if _, exists := prevFilteredObjects[id]; !exists {
				changes["added"].(map[string]Object)[id] = object
			}
		}

		for id := range prevFilteredObjects {
			if _, exists := filteredObjects[id]; !exists {
				changes["removed"] = append(changes["removed"].([]string), id)
			}
		}

		prevFilteredObjects = filteredObjects

		if err = conn.WriteJSON(changes); err != nil {
			fmt.Println("Error writing JSON:", err)
			break
		}
	}
}
