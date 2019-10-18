package voronoi

import (
	"fmt"
	pq "github.com/jupp0r/go-priority-queue"
	"math"
)

type Box struct {
	X0, X1, Y0, Y1 float64
}

type Voronoi struct {
	Result  []*Edge
	Arc     Arc
	Sites   pq.PriorityQueue
	Circles pq.PriorityQueue
	Box     Box
}

func NewVoronoi(points []Point) Voronoi {
	voronoi := Voronoi{}
	voronoi.Sites = CreateQ()
	voronoi.Circles = CreateQ()

	// bounding box
	voronoi.Box = Box{-50.0, -50.0, 550.0, 550.0}

	// insert points in the sites Q
	for _, el := range points {
		point := Point{el.X, el.Y}
		PushPoint(point, &voronoi.Sites)
		//Resize bounding box if there are points larger than its boundaries
		if point.X < voronoi.Box.X0 {
			voronoi.Box.X0 = point.X
		}
		if point.Y < voronoi.Box.Y0 {
			voronoi.Box.Y0 = point.Y
		}
		if point.X > voronoi.Box.X1 {
			voronoi.Box.X1 = point.X
		}
		if point.Y > voronoi.Box.Y1 {
			voronoi.Box.Y1 = point.Y
		}
	}

	// add margin to the bounding box
	dx := (voronoi.Box.X1 - voronoi.Box.X0 + 1) / 5.0
	dy := (voronoi.Box.Y1 - voronoi.Box.Y0 + 1) / 5.0
	voronoi.Box.X0 -= dx
	voronoi.Box.X1 += dx
	voronoi.Box.Y0 -= dy
	voronoi.Box.Y1 += dy

	return voronoi
}

func (vor *Voronoi) Process() {
	for !IsEmpty(&vor.Sites) {
		// Check whether the following event is circle or site event
		if !IsEmpty(&vor.Circles) && (TopOfEvents(&vor.Circles).X <= TopOfPoints(&vor.Sites).X) {
			vor.processCircles()
		} else {
			vor.processSites()
		}
	}
	// The remaining points are circles
	for !IsEmpty(&vor.Circles) {
		vor.processCircles()
	}

	vor.Arc = vor.finishEdges()
	vor.printOutput()
}

func (vor *Voronoi) processCircles() {
	// Get next event from circle pq
	event := PopEvent(&vor.Circles)

	if event.Valid {
		// Start new edge
		edge := NewEdge(event.Point)
		vor.Result = append(vor.Result, &edge)

		// Remove associated arc (parabola)
		arc := event.Arc
		if arc.Previous != nil {
			arc.Previous.Next = arc.Next
			arc.Previous.RightEdge = &edge
		}
		if arc.Next != nil {
			arc.Next.Previous = arc.Previous
			arc.Next.LeftEdge = &edge
		}

		// Finish the edges before and after a
		if arc.LeftEdge != nil {
			arc.LeftEdge = Finish(event.Point, arc.LeftEdge)
		}
		if arc.RightEdge != nil {
			arc.RightEdge = Finish(event.Point, arc.RightEdge)
		}

		// Recheck circle events on either side of p
		if arc.Previous != nil {
			vor.checkCircleEvent(arc.Previous, event.X, event)
		}
		if arc.Next != nil {
			vor.checkCircleEvent(arc.Next, event.X, event)
		}

	}
}

func (vor *Voronoi) processSites() {
	// Get the next event from the sites queue
	point := PopPoint(&vor.Sites)

	// Since new point defines new parabola(arc), we create one and then we insert it.
	vor.insertArc(point)
}

func (vor *Voronoi) insertArc(point Point) {
	// First we check if we have an arc in our BeachLine (if the arc we are inserting is the first arc)
	// If there is no Arc in our Voronoi structure, we add it
	if (vor.Arc) == (Arc{}) {
		// Since the new arc is the ONLY arc, we don't add any neighbors
		vor.Arc = NewArc(point, nil, nil)
		return
	}

	// We loop till we find the arc that is behind our Point
	currentArc := &vor.Arc
	for *currentArc != (Arc{}) {
		okay, intersectionPoint := intersect(point, *currentArc)
		// check whether we have intersectionPoint
		if okay {
			nextArc := currentArc.Next
			if nextArc == nil {
				nextArc = &Arc{}
			}
			hasIntersection, _ := intersect(point, *nextArc)

			if (currentArc.Next != nil) && !hasIntersection {
				// Add arc between two existing arcs
				arc := Arc{
					Point:     currentArc.Point,
					Previous:  currentArc,
					Next:      currentArc.Next,
					LeftEdge:  nil,
					RightEdge: nil,
				}
				currentArc.Next.Previous = &arc
				currentArc.Next = arc.Next.Previous
			} else {
				arc := NewArc(currentArc.Point, currentArc, nil)
				currentArc.Next = &arc
			}
			// Add point between current arc and the next one
			currentArc.Next.RightEdge = currentArc.RightEdge

			newArc := NewArc(point, currentArc, currentArc.Next)
			currentArc.Next.Previous = &newArc
			currentArc.Next = currentArc.Next.Previous

			currentArc = currentArc.Next

			// Add new edge connected to the intersection point (we will expand this edge like a ray till we dont find another intersection point
			segment := NewEdge(intersectionPoint)
			vor.Result = append(vor.Result, &segment)
			currentArc.Previous.RightEdge, currentArc.LeftEdge = &segment, &segment

			// Add new edge connected to the intersection point (we will expand this edge like a ray till we dont find another intersection point
			segment = NewEdge(intersectionPoint)
			vor.Result = append(vor.Result, &segment)
			currentArc.Next.LeftEdge, currentArc.RightEdge = &segment, &segment

			// check whether the newly added event caused a circle event
			vor.checkCircleEvent(currentArc, point.X, currentArc.Event)
			vor.checkCircleEvent(currentArc.Previous, point.X, currentArc.Event)
			vor.checkCircleEvent(currentArc.Next, point.X, currentArc.Event)
			return
		}
		currentArc = currentArc.Next
	}
	// If the parabola defined by the point never intersects the arc
	// just add it to the list
	currentArc = &vor.Arc
	// the arc structure acts like a linked list, so we loop till we find the last arc
	for currentArc.Next != nil {
		currentArc = currentArc.Next
	}
	arc := NewArc(point, currentArc, nil)
	currentArc.Next = &arc

	// Insert new segment between the point and the arc
	x := vor.Box.X0
	y := (currentArc.Next.Point.Y + currentArc.Point.Y) / 2.0
	start := Point{x, y}
	segment := NewEdge(start)
	currentArc.RightEdge, currentArc.Next.LeftEdge = &segment, &segment
	vor.Result = append(vor.Result, &segment)
}

func (vor *Voronoi) checkCircleEvent(arc *Arc, x0 float64, eventD *Event) {
	// Check whether we have new circle event for the arc
	if arc.Event != nil && arc.Event.X != x0 {
		(*arc).Event.Valid = false
	}
	//arc.Event = nil
	if (arc.Previous == nil) || (arc.Next == nil) {
		return
	}

	hasCircle, x, center := vor.circle(arc.Previous.Point, arc.Point, arc.Next.Point)
	if hasCircle && x > vor.Box.X0 {
		event := Event{
			X:     x,
			Point: center,
			Arc:   arc,
			Valid: true,
		}
		arc.Event = &event
		PushEvent(&event, &vor.Circles)
	}
}

func (vor *Voronoi) circle(a Point, b Point, c Point) (bool, float64, Point) {
	// Check if BC makes 'right turn' from AB
	if ((b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y)) > 0 {
		return false, 0.0, Point{}
	}

	// Algorithm of Rourke
	A := b.X - a.X
	B := b.Y - a.Y
	C := c.X - a.X
	D := c.Y - a.Y
	E := A*(a.X+b.X) + B*(a.Y+b.Y)
	F := C*(a.X+c.X) + D*(a.Y+c.Y)
	G := 2 * (A*(c.Y-b.Y) - B*(c.X-b.X))

	// If G == 0 -> points are colinear
	if G == 0 {
		return false, 0.0, Point{}
	}

	// Point o is the center of the circle
	ox := 1.0 * (D*E - B*F) / G
	oy := 1.0 * (A*F - C*E) / G

	// o.x plus radius equals max x coord
	x := ox + math.Sqrt(math.Pow(a.X-ox, 2)+math.Pow(a.Y-oy, 2))
	o := Point{ox, oy}
	return true, x, o
}

/*
	Check whether a new parabola defined by point P intersect with Arc..
*/
func intersect(point Point, arc Arc) (bool, Point) {
	// Make sure we have an Arc before we try to find any intersection
	if arc == (Arc{}) {
		return false, Point{}
	}
	// if they are collinear/neighbours, we cannot yet determine the intersection point
	if arc.Point.X == point.X {
		return false, Point{}
	}

	a, b := 0.0, 0.0

	// First we check whether there is intersection with the neighbours
	if arc.Previous != nil {
		// The 3rd param defines a line
		a = intersection(arc.Previous.Point, arc.Point, 1.0*point.X).Y
	}
	if arc.Next != nil {
		// The 3rd param defines a line
		b = intersection(arc.Point, arc.Next.Point, 1.0*point.X).Y
	}

	// If there is no neighbor intersection or the arcs are too far apart to intersect
	// we try to determine when the two arc will intersect once they start expanding
	if ((arc.Previous == nil) || (a <= point.Y)) && ((arc.Next == nil) || (point.Y <= b)) {
		pY := point.Y
		// Parabola equation
		pX := 1.0 * (math.Pow(arc.Point.X, 2) + math.Pow(arc.Point.Y-pY, 2) - math.Pow(point.X, 2)) / (2*arc.Point.X - 2*point.X)
		return true, Point{pX, pY}
	}
	return false, Point{}
}

/**
Find the intersection point of two parabolas defined by two points
*/
func intersection(point0, point1 Point, line float64) Point {
	// Lets assume that point0 is our intersection point
	p := point0
	var pX, pY float64

	// If both points have equal X coordinate, they will intersect in the mid of their Y distance
	if point0.X == point1.X {
		pY = (point0.Y + point1.Y) / 2
	} else if point1.X == line {
		// if point1 lies on the line sweeper
		pY = point1.Y
	} else if point0.X == line {
		pY = point0.Y
		p = point1
	} else {
		// Use quadratic formula for parabola definition
		z0 := 2.0 * (point0.X - line)
		z1 := 2.0 * (point1.X - line)

		a := 1.0/z0 - 1.0/z1
		b := -2.0 * (point0.Y/z0 - point1.Y/z1)
		c := 1.0*(math.Pow(point0.Y, 2)+math.Pow(point0.X, 2)-math.Pow(line, 2))/z0 - 1.0*(math.Pow(point1.Y, 2)+math.Pow(point1.X, 2)-math.Pow(line, 2))/z1

		pY = 1.0 * (-b - math.Sqrt(b*b-4*a*c)) / (2 * a)
	}
	// Determine X based on the Y
	pX = 1.0 * ((p.X * p.X) + (p.Y-pY)*(p.Y-pY) - line*line) / (2*p.X - 2*line)
	return Point{pX, pY}
}

func (vor *Voronoi) printOutput() {
	for _, edge := range vor.Result {
		p0, p1 := (*edge).Start, (*edge).End
		fmt.Printf("%f, %f, %f, %f ", p0.X, p0.Y, p1.X, p1.Y)
		println()
	}
}

func (vor *Voronoi) finishEdges() Arc {
	line := vor.Box.X1 + (vor.Box.X1 - vor.Box.X0) + (vor.Box.Y1 - vor.Box.Y0)
	arc := &(*vor).Arc
	for arc.Next != nil {
		if arc.RightEdge != (Edge{}) {
			intersectPoint := intersection(arc.Point, arc.Next.Point, line*2.0)
			(*arc).RightEdge = Finish(intersectPoint, &(arc).RightEdge)
		}
		arc = arc.Next
	}
	return *arc
}

//func (vor *Voronoi) getOutput() []Point {
//	//for _, edge := range vor.Result {
//	//	counter++
//	//	p0, p1 := edge.Start, edge.End
//	//	fmt.Printf("%f, %f, %f, %f", p0.X, p0.Y, p1.X, p1.Y)
//	//}
//}
//
