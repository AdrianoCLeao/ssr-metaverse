package server

type World struct {
	Objects map[string]Object
}

type Object struct {
	ID       string    `json:"id"`
	ModelURL string    `json:"model_url"`
	Position [3]float64 `json:"position"`
	Rotation [3]float64 `json:"rotation"`
}

func NewWorld() *World {
    return &World{
        Objects: map[string]Object{
            "world": {
                ID:       "world",
                ModelURL: "http://localhost:8080/assets/models/world.glb",
                Position: [3]float64{0, 0, 0},
                Rotation: [3]float64{0, 0, 0},
            },
            "navmesh": {
                ID:       "navmesh",
                ModelURL: "http://localhost:8080/assets/models/newmesh.gltf",
                Position: [3]float64{0, 0, 0},
                Rotation: [3]float64{0, 0, 0},
            },
        },
    }
}
