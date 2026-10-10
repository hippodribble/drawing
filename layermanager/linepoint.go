package layermanager

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

// Polyline represents a sequence of connected points in 2D space. It is defined by a slice of Points.
type Polyline struct {
	Points []Point
}

// NewPolyline creates a new Polyline with the given slice of Points.
func NewPolyline(points []Point) Polyline {
	return Polyline{
		Points: points,
	}
}