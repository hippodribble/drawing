package layermanager

import "math"

// Point represents a point in 2D space with X and Y coordinates.
type Point struct {
	X, Y float64
}

// NewPoint creates a new Point with the given X and Y coordinates.
func NewPoint(x1, y1 float64) Point {
	return Point{
		X: x1,
		Y: y1,
	}
}

// SquareDistance calculates the squared distance between the current Point and another Point.
func (p Point) SquareDistance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return dx*dx + dy*dy
}

// Line represents a line segment in 2D space defined by two endpoints (X1, Y1) and (X2, Y2).
type Line struct {
	X1, Y1 float64
	X2, Y2 float64
}

// NewLine creates a new Line with the given endpoints (X1, Y1) and (X2, Y2).
func NewLine(x1, y1, x2, y2 float64) Line {
	return Line{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}
}

// Length calculates the length of the line segment using the distance formula.
func (l Line) Length() float64 {
	dx := l.X2 - l.X1
	dy := l.Y2 - l.Y1
	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceTo calculates the shortest distance from the Line segment to a given Point. It uses the projection of the point onto the line segment to determine the closest point on the segment and then calculates the distance to that point.
func (l Line) DistanceTo(p Point) float64 {
	// Calculate the squared length of the line segment
	lineLengthSquared := (l.X2-l.X1)*(l.X2-l.X1) + (l.Y2-l.Y1)*(l.Y2-l.Y1)
	if lineLengthSquared == 0 {
		// The line segment is a point, return distance to that point
		return math.Sqrt((p.X-l.X1)*(p.X-l.X1) + (p.Y-l.Y1)*(p.Y-l.Y1))
	}

	// Calculate the projection of point p onto the line segment
	t := ((p.X-l.X1)*(l.X2-l.X1) + (p.Y-l.Y1)*(l.Y2-l.Y1)) / lineLengthSquared
	t = math.Max(0, math.Min(1, t)) // Clamp t to the range [0, 1]

	// Find the closest point on the line segment
	closestX := l.X1 + t*(l.X2-l.X1)
	closestY := l.Y1 + t*(l.Y2-l.Y1)

	// Calculate the distance from point p to the closest point on the line segment
	dx := p.X - closestX
	dy := p.Y - closestY
	return math.Sqrt(dx*dx + dy*dy)
}

// Polyline represents a sequence of connected points in 2D space. It is defined by a slice of Points.
type Polyline struct {
	ID string
	Points []Point
}

// NewPolyline creates a new Polyline with the given slice of Points.
func NewPolyline(points []Point) Polyline {
	return Polyline{
		Points: points,
	}
}

func (p Polyline) Length() float64 {
	if len(p.Points) < 2 {
		return 0
	}
	length := 0.0
	for i := 1; i < len(p.Points); i++ {
		dx := p.Points[i].X - p.Points[i-1].X
		dy := p.Points[i].Y - p.Points[i-1].Y
		length += math.Sqrt(dx*dx + dy*dy)
	}
	return length
}
