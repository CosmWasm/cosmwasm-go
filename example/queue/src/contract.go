package src

type Enqueue struct {
	Value int32 `json:"value"`
}

type Dequeue struct{}

type ExecuteMsg struct {
	// Enqueue adds a value in the queue
	Enqueue *Enqueue `json:"enqueue"`
	// Dequeue removes a value from the queue
	Dequeue *Dequeue `json:"dequeue"`
}

type QueryMsg struct {
	// Count counts how many items in the queue
	Count *struct{} `json:"count"`
	// Sum the number of values in the queue
	Sum *struct{} `json:"sum"`
	// Reducer keeps open two iters at once
	Reducer *struct{} `json:"reducer"`
	// List lists
	List *struct{} `json:"list"`
}
