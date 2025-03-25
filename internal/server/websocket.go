package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"ssr-metaverse/internal/database"
	"sync"

	"github.com/gorilla/websocket"
)

/*
   Represents a connected client.
   ID: A string that uniquely identifies the client (for example, client's remote address).
   Conn: The WebSocket connection associated with this client.
*/
type Client struct {
	ID   string
	Conn *websocket.Conn
}

/*
   Server holds:
   - A map of connected clients.
   - A mutex for thread-safe access to the Clients map.
   - A WebSocket upgrader to handle the handshake from HTTP to WebSocket.
*/
type Server struct {
	Clients  map[string]*Client
	Mutex    sync.Mutex
	Upgrader websocket.Upgrader
	DB       database.DBInterface
	Minio    database.MinioInterface
	Mongo 	 database.MongoInterface
	Redis 	 database.RedisInterface
}

/*
   isWithinDistance calculates the Euclidean distance between the client and
   the object, then returns true if the distance is less than or equal
   to maxDistance.
*/
func isWithinDistance(clientPos, objectPos [3]float64, maxDistance float64) bool {
	dx := clientPos[0] - objectPos[0]
	dy := clientPos[1] - objectPos[1]
	dz := clientPos[2] - objectPos[2]
	distance := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return distance <= maxDistance
}

/*
   filterObjects iterates through all objects in the world and returns only
   those that are within maxDistance of the client's position.
*/
func filterObjects(world *World, clientPos [3]float64, maxDistance float64) map[string]Object {
	filtered := make(map[string]Object)
	for id, object := range world.Objects {
		if isWithinDistance(clientPos, object.Position, maxDistance) {
			filtered[id] = object
		}
	}
	return filtered
}

/*
   Creates and initializes a new Server instance.
*/
func NewServer(db database.DBInterface, minio database.MinioInterface, mongo database.MongoInterface, redis database.RedisInterface) *Server{
	return &Server{
		Clients: make(map[string]*Client),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		DB:    db,
		Minio: minio, 
		Mongo: mongo,
		Redis: redis,
	}
}

/*
   HandleWebSocket is the main WebSocket logic:
   - Upgrades the connection from HTTP to WebSocket.
   - Adds the client to the server's map of connected clients.
   - Listens for messages (position updates) from the client.
   - Filters world objects based on the client position, sends back "diffs."
*/
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

	/*
	   Lock the mutex to safely add the client to the map,
	   then unlock before proceeding.
	*/
	s.Mutex.Lock()
	s.Clients[client.ID] = client
	s.Mutex.Unlock()

	/*
	   When the function returns, remove the client from the map
	   to clean up resources.
	*/
	defer func() {
		s.Mutex.Lock()
		delete(s.Clients, client.ID)
		s.Mutex.Unlock()
	}()

	/*
	   Create an example world. In a real application, this would
	   likely be a shared resource or loaded from elsewhere.
	*/
	world := NewWorld()

	/*
	   Keep track of the previously filtered objects for the purpose
	   of determining which objects are added or removed each tick.
	*/
	prevFilteredObjects := make(map[string]Object)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		/*
		   The payload should contain position and rotation data
		   from the client, encoded in JSON.
		*/
		var payload struct {
			Position [3]float64 `json:"position"`
			Rotation [3]float64 `json:"rotation"`
		}
		if err := json.Unmarshal(msg, &payload); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		/*
		   Filter the objects by specifying a maximum distance
		   the client can see.
		*/
		maxDistance := 50.0
		filteredObjects := filterObjects(world, payload.Position, maxDistance)

		/*
		   Prepare a structure to capture which objects are added
		   and which are removed since the last iteration.
		*/
		changes := map[string]interface{}{
			"added":   make(map[string]Object),
			"removed": make([]string, 0),
		}

		/*
		   Detect newly added objects.
		*/
		for id, object := range filteredObjects {
			if _, exists := prevFilteredObjects[id]; !exists {
				changes["added"].(map[string]Object)[id] = object
			}
		}

		/*
		   Detect removed objects that were previously visible but
		   are no longer in range.
		*/
		for id := range prevFilteredObjects {
			if _, exists := filteredObjects[id]; !exists {
				changes["removed"] = append(changes["removed"].([]string), id)
			}
		}

		/*
		   Update the previous filtered object set for the next iteration.
		*/
		prevFilteredObjects = filteredObjects

		/*
		   Write the changes back to the client.
		*/
		if err = conn.WriteJSON(changes); err != nil {
			fmt.Println("Error writing JSON:", err)
			break
		}
	}
}
