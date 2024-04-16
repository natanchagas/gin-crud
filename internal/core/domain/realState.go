package domain

type RealState struct {
	Id           uint64  `json:"id,omitempty"`
	Registration uint64  `json:"registration"`
	Address      string  `json:"address"`
	Size         uint64  `json:"size"`
	Price        float64 `json:"price"`
	State        string  `json:"state"`
}
