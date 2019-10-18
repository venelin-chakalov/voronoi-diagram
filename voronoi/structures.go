package voronoi

import "github.com/jupp0r/go-priority-queue"

func CreateQ() pq.PriorityQueue {
	return pq.New()
}

func PushPoint(point Point, queue *pq.PriorityQueue) {
	(*queue).Insert(point, point.X)
}

func PushEvent(event *Event, queue *pq.PriorityQueue) {
	(*queue).Insert(event, event.X)
}

func PopPoint(queue *pq.PriorityQueue) Point {
	lastEl, _ := queue.Pop()
	return lastEl.(Point)
}
func PopEvent(queue *pq.PriorityQueue) *Event {
	lastEl, _ := (*queue).Pop()
	return lastEl.(*Event)
}

func IsEmpty(queue *pq.PriorityQueue) bool {
	return queue.Len() == 0
}

//TODO REMOVE CASTING
func TopOfPoints(queue *pq.PriorityQueue) Point {
	point, _ := (*queue).Pop()
	(*queue).Insert(point, point.(Point).X)
	return point.(Point)
}

func TopOfEvents(queue *pq.PriorityQueue) *Event {
	event := PopEvent(queue)
	PushEvent(event, queue)
	return event
}

type Point struct {
	X, Y float64
}

type Event struct {
	X     float64
	Point Point
	Arc   *Arc
	Valid bool
}

type Arc struct {
	Point     Point
	Previous  *Arc
	Next      *Arc
	LeftEdge  Edge
	RightEdge Edge
	Event     *Event
}

func NewArc(point Point, previous, next *Arc) Arc {
	return Arc{
		Point:     point,
		Previous:  previous,
		Next:      next,
		LeftEdge:  Edge{},
		RightEdge: Edge{},
		Event:     nil,
	}
}

type Edge struct {
	Start Point
	End   Point
	Done  bool
}

func NewEdge(point Point) Edge {
	return Edge{
		Start: point,
		End:   Point{},
		Done:  false,
	}
}

func Finish(point Point, edge Edge) Edge {
	if (edge).Done {
		return edge
	}
	(edge).End = point
	(edge).Done = true
	return edge
}
