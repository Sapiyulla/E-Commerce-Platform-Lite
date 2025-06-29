package dto

type Status struct {
	Code    int    `json:"-"` //0 - Not | 1 - Client | 2 - Server
	Status  string `json:"status"`
	Message string `json:"message"`
}

// statuses
var (
	Error    = "error"
	Success  = "success"
	NotFound = "not found"
)
