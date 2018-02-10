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
}

type PointsContainer struct{
	Point_list []Point `json:"points"`
	Point_metadata Point_metadata `json:"metadata"`
}

type Graph_metadata struct{
	ID string `json:"id"`
	Column_x string `json:"column_x"`
	Column_y string `json:"column_y"`
}

type Graph struct{
	Graph_metadata Graph_metadata `json:"metadata"`
	PointsContainer []PointsContainer `json:"points_container"`
}

type HttpData []Graph