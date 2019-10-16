package main

import (
	"voronoi"
)

func main() {
	voroni := voronoi.NewVoronoi([]voronoi.Point{
		{89.0, 379.0},
		{198.0, 179.0},
		{428.0, 268.0},
	})
	voroni.Process()
	println(voroni.Result)
}
