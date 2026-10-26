package task

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type CreateTaskRequest struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}
