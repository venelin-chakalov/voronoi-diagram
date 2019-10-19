package test

import (
	"fmt"
	"voronoi"
)

func main() {
	voronoi := voronoi.NewVoronoi([]voronoi.Point{
		{89.0, 379.0},
		{198.0, 179.0},
		{428.0, 268.0},
	})
	voronoi.Process()
	result := voronoi.Result
	for _, edge := range result {
		fmt.Printf("%v '\n", edge)
	}
}
