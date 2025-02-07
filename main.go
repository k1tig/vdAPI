package main

type racer struct {
	Name string `json:"racername"`
}

var racers []racer

func main() {
	server := NewAPIServer((":8080"))
	server.Run()
}
