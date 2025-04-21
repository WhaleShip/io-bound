package domain

type TaskStatus struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}
