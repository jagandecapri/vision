package server

type Point_data struct{
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Point_metadata struct{
	Color string `json:"color"`
}

type Point struct{
	Point_data Point_data `json:"data"`
	Point_metadata Point_metadata `json:"metadata"`
}

type Graph_metadata struct{
	ID string `json:"id"`
}

type Graph struct{
	Graph_metadata Graph_metadata `json:"metadata"`
	Points []Point `json:"points"`
}

type HttpData []Graph