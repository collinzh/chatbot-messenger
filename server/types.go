package server

type ChatRequest struct {
	Input string `json:"input"`
}

type ChatResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}
