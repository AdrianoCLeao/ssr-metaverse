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
            "object1": {
                ID:       "object1",
                ModelURL: "http://localhost:8080/assets/models/world.glb",
                Position: [3]float64{0, 0, 0},
                Rotation: [3]float64{0, 0, 0},
            },
            "object2": {
                ID:       "object2",
                ModelURL: "http://localhost:8080/assets/avatars/user.glb",
                Position: [3]float64{1, 0, 1},
                Rotation: [3]float64{0, 45, 0},
            },
        },
    }
}
