package dto

// Example represents a resource in the API.
type Example struct {
	ID   string `json:"id" doc:"Resource ID"`
	Name string `json:"name" doc:"Resource name"`
}
