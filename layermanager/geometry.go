package layermanager

import (
	"math"
)

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

// Translate moves the point by the supplied dx and dy offset.
func (p Point) Translate(dx, dy float64) Point {
	p.X += dx
	p.Y += dy
	return p
}

// Rotate moves the point around the supplied centre by the given angle in radians.
func (p Point) Rotate(angle float64, centre Point) Point {
	dx := p.X - centre.X
	dy := p.Y - centre.Y
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)
	return Point{
		X: centre.X + dx*cosTheta - dy*sinTheta,
		Y: centre.Y + dx*sinTheta + dy*cosTheta,
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

// Translate moves both endpoints of the line by the supplied dx and dy offset.
func (l Line) Translate(dx, dy float64) Line {
	l.X1 += dx
	l.Y1 += dy
	l.X2 += dx
	l.Y2 += dy
	return l
}

// Rotate moves both endpoints of the line around the supplied centre by the given angle in radians.
func (l Line) Rotate(angle float64, centre Point) Line {
	p1 := Point{X: l.X1, Y: l.Y1}.Rotate(angle, centre)
	p2 := Point{X: l.X2, Y: l.Y2}.Rotate(angle, centre)
	return Line{
		X1: p1.X,
		Y1: p1.Y,
		X2: p2.X,
		Y2: p2.Y,
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

// Polygon represents a closed shape defined by a sequence of connected points in 2D space. It is defined by a slice of Points.
type Polygon struct {
	ID     string
	Points []Point
}

// NewPolygon creates a new Polygon with the given slice of Points.
func NewPolygon(points []Point) Polygon {
	return NewPolyline(points).Close()
}

// Polyline represents a sequence of connected points in 2D space. It is defined by a slice of Points.
type Polyline struct {
	ID     string
	Points []Point
}

// NewPolyline creates a new Polyline with the given slice of Points.
func NewPolyline(points []Point) Polyline {
	return Polyline{
		Points: points,
	}
}

// Translate moves each point in the polyline by the supplied dx and dy offset.
func (p Polyline) Translate(dx, dy float64) Polyline {
	translated := Polyline{ID: p.ID, Points: make([]Point, len(p.Points))}
	for i, point := range p.Points {
		translated.Points[i] = point.Translate(dx, dy)
	}
	return translated
}

// Rotate moves each point in the polyline around the supplied centre by the given angle in radians.
func (p Polyline) Rotate(angle float64, centre Point) Polyline {
	rotated := Polyline{ID: p.ID, Points: make([]Point, len(p.Points))}
	for i, point := range p.Points {
		rotated.Points[i] = point.Rotate(angle, centre)
	}
	return rotated
}

// Length calculates the total length of the polyline by summing the lengths of its constituent line segments.
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

// DouglasPuecker returns a simplified approximation of the polyline using
// the Douglas-Peucker algorithm and the supplied accuracy as a tolerance.
func (p Polyline) DouglasPuecker(accuracy float64) Polyline {
	if len(p.Points) < 3 || accuracy <= 0 {
		points := make([]Point, len(p.Points))
		copy(points, p.Points)
		return Polyline{ID: p.ID, Points: points}
	}

	points := make([]Point, len(p.Points))
	copy(points, p.Points)

	closed := points[0].SquareDistance(points[len(points)-1]) <= .001
	if closed {
		points = points[:len(points)-1]
		if len(points) < 3 {
			out := make([]Point, len(p.Points))
			copy(out, p.Points)
			return Polyline{ID: p.ID, Points: out}
		}
	}

	simplified := douglasPueckerPoints(points, accuracy)
	if closed {
		simplified = append(simplified, simplified[0])
	}

	return Polyline{ID: p.ID, Points: simplified}
}

// Douglas-Puecker recursion
func douglasPueckerPoints(points []Point, tolerance float64) []Point {
	if len(points) <= 2 {
		out := make([]Point, len(points))
		copy(out, points)
		return out
	}

	start := points[0]
	end := points[len(points)-1]
	maxDistance := -1.0
	maxIndex := 0

	for i := 1; i < len(points)-1; i++ {
		distance := pointToSegmentDistance(points[i], start, end)
		if distance > maxDistance {
			maxDistance = distance
			maxIndex = i
		}
	}

	if maxDistance <= tolerance {
		return []Point{start, end}
	}

	left := douglasPueckerPoints(points[:maxIndex+1], tolerance)
	right := douglasPueckerPoints(points[maxIndex:], tolerance)

	out := make([]Point, 0, len(left)+len(right)-1)
	out = append(out, left[:len(left)-1]...)
	out = append(out, right...)
	return out
}

// pointToSegmentDistance calculates the shortest distance from a point to a line segment defined by two endpoints (start and end).
// It uses vector projection to find the closest point on the segment and then computes the distance to that point.
func pointToSegmentDistance(point, start, end Point) float64 {
	segmentLengthSquared := start.SquareDistance(end)
	if segmentLengthSquared == 0 {
		return math.Sqrt(point.SquareDistance(start))
	}

	t := ((point.X-start.X)*(end.X-start.X) + (point.Y-start.Y)*(end.Y-start.Y)) / segmentLengthSquared
	t = math.Max(0, math.Min(1, t))

	projection := Point{
		X: start.X + t*(end.X-start.X),
		Y: start.Y + t*(end.Y-start.Y),
	}
	return math.Sqrt(point.SquareDistance(projection))
}

// Close checks if the first and last points of the Polyline are the same. If they are not,
// it appends the first point to the end of the Points slice, effectively closing the Polyline.
// It returns the modified Polyline.
func (p Polyline) Close() Polygon {
	if len(p.Points) < 2 {
		return Polygon{ID: p.ID, Points: p.Points}
	}
	if p.Points[0].SquareDistance(p.Points[len(p.Points)-1]) > .001 {
		p.Points = append(p.Points, p.Points[0])
	}
	return Polygon{ID: p.ID, Points: p.Points}
}

// Close checks if the first and last points of the Polyline are the same. If they are not,
// it appends the first point to the end of the Points slice, effectively closing the Polygon.
// It returns the modified Polygon.
func (p Polygon) Close() Polygon {
	if len(p.Points) < 2 {
		return Polygon{ID: p.ID, Points: p.Points}
	}
	if p.Points[0].SquareDistance(p.Points[len(p.Points)-1]) > .001 {
		p.Points = append(p.Points, p.Points[0])
	}
	return Polygon{ID: p.ID, Points: p.Points}
}

// Translate moves each point in the polygon by the supplied dx and dy offset.
func (p Polygon) Translate(dx, dy float64) Polygon {
	translated := Polygon{ID: p.ID, Points: make([]Point, len(p.Points))}
	for i, point := range p.Points {
		translated.Points[i] = point.Translate(dx, dy)
	}
	return translated
}

// Rotate moves each point in the polygon around the supplied centre by the given angle in radians.
func (p Polygon) Rotate(angle float64, centre Point) Polygon {
	rotated := Polygon{ID: p.ID, Points: make([]Point, len(p.Points))}
	for i, point := range p.Points {
		rotated.Points[i] = point.Rotate(angle, centre)
	}
	return rotated
}

// Offset creates a new Polyline that is offset from the original Polyline by a specified distance.
// The offset is applied perpendicular to the direction of the Polyline at each point.
// A positive offset value moves the Polyline to the left, while a negative value moves it to the right.
func (p Polyline) Offset(offset float64) Polyline {
	if len(p.Points) < 2 {
		return p
	}

	offsetPoints := make([]Point, len(p.Points))
	for i := 0; i < len(p.Points); i++ {
		var dx, dy float64
		if i == 0 {
			dx = p.Points[i+1].X - p.Points[i].X
			dy = p.Points[i+1].Y - p.Points[i].Y
		} else if i == len(p.Points)-1 {
			dx = p.Points[i].X - p.Points[i-1].X
			dy = p.Points[i].Y - p.Points[i-1].Y
		} else {
			dx = p.Points[i+1].X - p.Points[i-1].X
			dy = p.Points[i+1].Y - p.Points[i-1].Y
		}

		length := math.Sqrt(dx*dx + dy*dy)
		if length == 0 {
			offsetPoints[i] = p.Points[i]
			continue
		}

		// Normalize the direction vector and calculate the perpendicular offset
		perpX := -dy / length
		perpY := dx / length

		offsetPoints[i] = Point{
			X: p.Points[i].X + perpX*offset,
			Y: p.Points[i].Y + perpY*offset,
		}
	}

	return Polyline{
		ID:     p.ID,
		Points: offsetPoints,
	}
}

// BufferBy creates a new Polygon that is offset from the original Polygon by a specified distance.
func (q Polygon) offset(offset float64) Polygon {
	if len(q.Points) < 2 {
		return q
	}
	// N := len(q.Points)

	qp := q.Points
	qp = append(qp, qp[1], qp[2])

	offsetPoints := make([]Point, len(qp))

	for i := 0; i < len(qp); i++ {
		prevIndex := (i - 1 + len(qp)) % len(qp)
		nextIndex := (i + 1) % len(qp)

		// Edge vectors
		edge1 := Point{X: qp[i].X - qp[prevIndex].X, Y: qp[i].Y - qp[prevIndex].Y}
		edge2 := Point{X: qp[nextIndex].X - qp[i].X, Y: qp[nextIndex].Y - qp[i].Y}

		// Compute edge normals
		normal1 := Point{X: -edge1.Y, Y: edge1.X}
		normal2 := Point{X: -edge2.Y, Y: edge2.X}

		// Normalize the normals
		norm1Length := math.Sqrt(normal1.X*normal1.X + normal1.Y*normal1.Y)
		norm2Length := math.Sqrt(normal2.X*normal2.X + normal2.Y*normal2.Y)
		if norm1Length > 0 {
			normal1.X /= norm1Length
			normal1.Y /= norm1Length
		}
		if norm2Length > 0 {
			normal2.X /= norm2Length
			normal2.Y /= norm2Length
		}

		// Calculate average of normals for the offset
		averageNormal := Point{
			X: (normal1.X + normal2.X) / 2,
			Y: (normal1.Y + normal2.Y) / 2,
		}

		// Normalize the average normal
		averageNormalLength := math.Sqrt(averageNormal.X*averageNormal.X + averageNormal.Y*averageNormal.Y)
		if averageNormalLength > 0 {
			averageNormal.X /= averageNormalLength
			averageNormal.Y /= averageNormalLength
		}

		// Create the offset point
		offsetPoints[i] = Point{
			X: qp[i].X + averageNormal.X*offset,
			Y: qp[i].Y + averageNormal.Y*offset,
		}
	}

	return Polygon{
		ID:     q.ID,
		Points: offsetPoints[1 : len(offsetPoints)-1],
	}
}

// Area calculates the area of the Polygon using the shoelace formula. It iterates through the points of the polygon,
// summing the cross products of adjacent points to compute the area. The result is divided by 2 to obtain the final area value.
// It is a signed area, so the result will be negative if the points are ordered clockwise. The absolute value of the area can be taken if needed.
func (p Polygon) Area() float64 {
	area := 0.0
	n := len(p.Points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += p.Points[i].X * p.Points[j].Y
		area -= p.Points[j].X * p.Points[i].Y
	}
	return area / 2.0
}

// Perimeter calculates the perimeter of the Polygon by summing the distances between consecutive points.
func (p Polygon) Perimeter() float64 {
	perimeter := 0.0
	n := len(p.Points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		dx := p.Points[j].X - p.Points[i].X
		dy := p.Points[j].Y - p.Points[i].Y
		perimeter += math.Sqrt(dx*dx + dy*dy)
	}
	return perimeter
}

// Centroid calculates the centroid (geometric center) of the Polygon.
func (p Polygon) Centroid() Point {
	centroid := Point{X: 0, Y: 0}
	area := p.Area()
	if area == 0 {
		return centroid
	}

	n := len(p.Points)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		factor := (p.Points[i].X*p.Points[j].Y - p.Points[j].X*p.Points[i].Y)
		centroid.X += (p.Points[i].X + p.Points[j].X) * factor
		centroid.Y += (p.Points[i].Y + p.Points[j].Y) * factor
	}
	centroid.X /= (6 * area)
	centroid.Y /= (6 * area)

	return centroid
}

// IsClockwise determines if the points of the Polygon are ordered in a clockwise direction.
// It uses the area calculation to determine the orientation of the points. If the area is negative,
// the points are ordered clockwise; if positive, they are ordered counter-clockwise.
func (p Polygon) IsClockwise() bool {
	return p.Area() < 0
}

// Reverse reverses the order of the points in the Polygon. This can be useful for changing the
// orientation of the polygon from clockwise to counter-clockwise or vice versa.
func (p Polygon) Reverse() Polygon {
	q := Polygon{ID: p.ID, Points: make([]Point, len(p.Points))}
	copy(q.Points, p.Points)
	n := len(p.Points)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		p.Points[i], p.Points[j] = p.Points[j], p.Points[i]
	}
	return q
}

// BufferBy creates a new Polygon that is offset from the original Polygon by a specified distance.
// The offset is applied perpendicular to the edges of the Polygon. A positive offset value moves the Polygon outward,
// while a negative value moves it inward. If the Polygon is ordered clockwise, it will be reversed before applying the offset.
func (p Polygon) BufferBy(r float64) Polygon {
	if p.IsClockwise() {
		p.Reverse()
	}
	return p.offset(-r)
}

// DouglasPeucker returns a simplified approximation of the polygon using
// the Douglas-Peucker algorithm and the supplied accuracy as a tolerance.
// Internally, it ccloses the polygon if necessary, converts it to a polyline,
// applies the Douglas-Peucker algorithm, and then closes the resulting polyline
// to form a polygon.
func (p Polygon) DouglasPeucker(tolerance float64) Polygon {
	points := p.Close().Points
	polyline := Polyline{Points: points, ID: p.ID}
	return polyline.DouglasPuecker(tolerance).Close()
}
