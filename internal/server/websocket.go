package server

import (
	"fmt"
	"net/http"
	"sync"
	"math"

	"github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
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

    clientPos := [3]float64{0, 0, 0}

    world := NewWorld()
    filteredObjects := filterObjects(world, clientPos, 5.0) 
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