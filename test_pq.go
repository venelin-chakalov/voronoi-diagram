package main

import "voronoi"

func main() {
	points := voronoi.CreateQ()
	events := voronoi.CreateQ()

	point := voronoi.Point{
		X: 1,
		Y: 2,
	}
	point2 := voronoi.Point{
		X: 2,
		Y: 3,
	}
	point3 := voronoi.Point{
		X: 3,
		Y: 4,
	}

	event := voronoi.Event{
		X:     1,
		Point: voronoi.Point{1,2},
		Arc:   voronoi.Arc{},
		Valid: false,
	}
	event2 := voronoi.Event{
		X:     2,
		Point: voronoi.Point{2,2},
		Arc:   voronoi.Arc{},
		Valid: false,
	}

	event3 := voronoi.Event{
		X:     3,
		Point: voronoi.Point{3,2},
		Arc:   voronoi.Arc{},
		Valid: false,
	}

	println("points")
	voronoi.PushPoint(point, &points)
	voronoi.PushPoint(point2, &points)
	voronoi.PushPoint(point3, &points)
	println(voronoi.TopOfPoints(&points).X)
	println(voronoi.PopPoint(&points).X)
	println(voronoi.TopOfPoints(&points).X)

	println("circles")
	voronoi.PushEvent(event, &events)
	voronoi.PushEvent(event2, &events)
	voronoi.PushEvent(event3, &events)
	println(voronoi.TopOfEvents(&events).X)
	println(voronoi.PopEvent(&events).X)
	println(voronoi.TopOfEvents(&events).X)
}
