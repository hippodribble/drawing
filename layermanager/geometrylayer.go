package layermanager

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GeometryLayer is a widget that represents a layer of geometric shapes (points and lines) that can be drawn on a canvas. It holds a collection of points and lines, along with a projection for converting between world and screen coordinates. The layer can be rendered with specified stroke width and color.
type GeometryLayer struct {
	widget.BaseWidget
	Points                 []Point
	Lines                  []Line
	Polylines              []Polyline
	Polygons               []Polygon
	Projection             *Projection
	statusChannel          chan string
	StrokeWidth            float32
	StrokeColor, fillColor color.Color
	nearestPoint           int
	nearestLine            int
	highlightLineStart     int
	highlightLineEnd       int
	highlightNearest       bool
}

// NewGeometryLayer creates a new GeometryLayer with the given projection. It initializes the layer with empty slices for points and lines, and sets the stroke width and color to default values.
func NewGeometryLayer() *GeometryLayer {
	layer := &GeometryLayer{
		Projection:         NewUnitProjection(),
		StrokeWidth:        1,
		StrokeColor:        theme.Color(theme.ColorNameForeground),
		nearestPoint:       -1,
		nearestLine:        -1,
		highlightLineStart: -1,
		highlightLineEnd:   -1,
	}
	layer.ExtendBaseWidget(layer)
	return layer
}

// CreateRenderer creates a new GeometryLayerRenderer for the GeometryLayer. The renderer is responsible for drawing the points and lines in the layer on the canvas, using the specified projection and stroke properties.
func (layer *GeometryLayer) CreateRenderer() fyne.WidgetRenderer {
	r := &GeometryLayerRenderer{
		layer: layer,
	}
	r.Sync()
	return r
}

// SetProjection sets the projection for the GeometryLayer and refreshes the layer to reflect the new projection. The projection is used to convert between world and screen coordinates when rendering the points and lines in the layer.
func (gl *GeometryLayer) SetProjection(proj *Projection) {
	gl.Projection = proj
	if gl.statusChannel != nil {
		gl.statusChannel <- "Projection changed to " + proj.String()
	}
	gl.invalidate()
}

// SetStatusChannel sets the status channel for the GeometryLayer. The status channel can be used to send status updates or messages related to the layer's operations, such as drawing or refreshing. This is generally used to notify the containing application of relevant changes.
func (gl *GeometryLayer) SetStatusChannel(ch chan string) {
	gl.statusChannel = ch
}

func (gl *GeometryLayer) SetHighlightNearest(b bool) {
	prev := gl.highlightNearest
	gl.highlightNearest = b
	if !b {
		gl.nearestPoint = -1
		gl.nearestLine = -1
		gl.highlightLineStart = -1
		gl.highlightLineEnd = -1
	}
	if prev != b {
		gl.invalidate()
	}
}

// SetHighlightedLineRange sets the contiguous line-index range that should be
// highlighted. Invalid ranges clear line highlighting.
func (gl *GeometryLayer) SetHighlightedLineRange(start, end int) {
	prevStart := gl.highlightLineStart
	prevEnd := gl.highlightLineEnd

	if start > end {
		start, end = end, start
	}

	if len(gl.Lines) == 0 || end < 0 || start >= len(gl.Lines) {
		gl.highlightLineStart = -1
		gl.highlightLineEnd = -1
	} else {
		if start < 0 {
			start = 0
		}
		if end >= len(gl.Lines) {
			end = len(gl.Lines) - 1
		}
		gl.highlightLineStart = start
		gl.highlightLineEnd = end
		gl.highlightNearest = true
		gl.nearestLine = start
	}

	if prevStart != gl.highlightLineStart || prevEnd != gl.highlightLineEnd {
		gl.invalidate()
	}
}

// HighlightPointByIndex highlights a specific point index in this layer.
// It enables highlighting automatically. Out-of-range values clear the highlight.
func (gl *GeometryLayer) HighlightPointByIndex(index int) {
	gl.highlightNearest = true
	if index < 0 || index >= len(gl.Points) {
		if gl.nearestPoint != -1 {
			gl.nearestPoint = -1
			gl.invalidate()
		}
		return
	}
	if gl.nearestPoint != index {
		gl.nearestPoint = index
		gl.invalidate()
	}
}

// ClearHighlight clears any highlighted point in this layer.
func (gl *GeometryLayer) ClearHighlight() {
	if gl.nearestPoint != -1 || gl.nearestLine != -1 || gl.highlightNearest {
		gl.nearestPoint = -1
		gl.nearestLine = -1
		gl.highlightLineStart = -1
		gl.highlightLineEnd = -1
		gl.highlightNearest = false
		gl.invalidate()
	}
}

// Limits calculates and returns the minimum and maximum X and Y coordinates of all points and lines in the GeometryLayer. It returns four float32 values: minX, minY, maxX, and maxY. If the layer has no points or lines, it returns default limits of (0, 0, 1, 1). This method is useful for determining the bounding box of the geometric shapes in the layer.
func (gl *GeometryLayer) Limits() (minX, minY, maxX, maxY float32) {
	if len(gl.Points) == 0 && len(gl.Lines) == 0 {
		return 0, 0, 1, 1
	}
	minX, maxX, minY, maxY = 1.0e+20, -1.0e+20, 1.0e+20, -1.0e+20
	if len(gl.Points) > 0 {
		for _, point := range gl.Points {
			if float32(point.X) < minX {
				minX = float32(point.X)
			}
			if float32(point.X) > maxX {
				maxX = float32(point.X)
			}
			if float32(point.Y) < minY {
				minY = float32(point.Y)
			}
			if float32(point.Y) > maxY {
				maxY = float32(point.Y)
			}
		}
	}
	if len(gl.Lines) > 0 {
		for _, line := range gl.Lines {
			if float32(line.X1) < minX {
				minX = float32(line.X1)
			}
			if float32(line.X2) < minX {
				minX = float32(line.X2)
			}
			if float32(line.Y1) < minY {
				minY = float32(line.Y1)
			}
			if float32(line.Y2) < minY {
				minY = float32(line.Y2)
			}
			if float32(line.X1) > maxX {
				maxX = float32(line.X1)
			}
			if float32(line.X2) > maxX {
				maxX = float32(line.X2)
			}
			if float32(line.Y1) > maxY {
				maxY = float32(line.Y1)
			}
			if float32(line.Y2) > maxY {
				maxY = float32(line.Y2)
			}
		}
	}

	return
}

// Clear removes all points and lines from the GeometryLayer and refreshes the layer. This effectively resets the layer to an empty state, allowing for new geometric shapes to be added without any previous data.
func (gl *GeometryLayer) Clear() {
	gl.Points = []Point{}
	gl.Lines = []Line{}
	gl.nearestPoint = -1
	gl.nearestLine = -1
	gl.highlightLineStart = -1
	gl.highlightLineEnd = -1
	gl.invalidate()
}

// sets the stroke width for the GeometryLayer and refreshes the layer to reflect the new stroke width. The stroke width determines how thick the lines and points will be drawn on the canvas.
//   - If you don't do this, the line won't be visible
func (gl *GeometryLayer) SetStrokeWidth(width float32) {
	gl.StrokeWidth = width
	gl.invalidate()
}

// sets the stroke color for the GeometryLayer and refreshes the layer to reflect the new stroke color. The stroke color determines the color used to draw the lines and points on the canvas.
// - If you don't do this, the line won't be visible
func (gl *GeometryLayer) SetStrokeColor(color color.Color) {
	gl.StrokeColor = color
	gl.invalidate()
}

// used when making points - it is the interior of the circle that represents points.
func (gl *GeometryLayer) SetFillColor(color color.Color) {
	gl.fillColor = color
	gl.invalidate()
}

func (g *GeometryLayer) AttachTo(m *LayerManager) {
	m.AddLayer(g)
}

// AddPolyline appends polyline to this layer and expands it into line segments
// used by the renderer.
func (gl *GeometryLayer) AddPolygon(polygon Polygon) {
	gl.Polygons = append(gl.Polygons, polygon)
	for i := 0; i < len(polygon.Points)-1; i++ {
		line := Line{
			X1: polygon.Points[i].X,
			Y1: polygon.Points[i].Y,
			X2: polygon.Points[i+1].X,
			Y2: polygon.Points[i+1].Y,
		}
		gl.Lines = append(gl.Lines, line)
	}
	gl.invalidate()
}

// AddPolyline appends polyline to this layer and expands it into line segments
// used by the renderer.
func (gl *GeometryLayer) AddPolyline(polyline Polyline) {
	gl.Polylines = append(gl.Polylines, polyline)
	for i := 0; i < len(polyline.Points)-1; i++ {
		line := Line{
			X1: polyline.Points[i].X,
			Y1: polyline.Points[i].Y,
			X2: polyline.Points[i+1].X,
			Y2: polyline.Points[i+1].Y,
		}
		gl.Lines = append(gl.Lines, line)
	}
	gl.invalidate()
}

// ClearPolylines removes all polylines and clears the generated line segments.
func (gl *GeometryLayer) ClearPolylines() {
	gl.Polylines = gl.Polylines[:0]
	gl.Lines = gl.Lines[:0]
	gl.highlightLineStart = -1
	gl.highlightLineEnd = -1
	gl.nearestLine = -1
	gl.invalidate()
}

// NearestPolyline returns the index of the polyline nearest to p and the
// minimum distance in world units.
func (gl *GeometryLayer) NearestPolyline(p Point) (nearestIndex int, distance float64) {
	nearestIndex = -1
	distance = math.Inf(1)

	for i, polyline := range gl.Polylines {
		if len(polyline.Points) < 2 {
			continue
		}

		minDistForPolyline := math.Inf(1)
		for j := 0; j < len(polyline.Points)-1; j++ {
			seg := Line{
				X1: polyline.Points[j].X,
				Y1: polyline.Points[j].Y,
				X2: polyline.Points[j+1].X,
				Y2: polyline.Points[j+1].Y,
			}
			d := seg.DistanceTo(p)
			if d < minDistForPolyline {
				minDistForPolyline = d
			}
		}

		if minDistForPolyline < distance {
			distance = minDistForPolyline
			nearestIndex = i
		}
	}

	return nearestIndex, distance
}

// HighlightNearestPolyline finds and highlights the nearest polyline to p.
func (gl *GeometryLayer) HighlightNearestPolyline(p Point) (nearestIndex int, distance float64) {
	nearestIndex, distance = gl.NearestPolyline(p)
	if nearestIndex < 0 {
		gl.SetHighlightedLineRange(-1, -1)
		return nearestIndex, distance
	}

	start, end, ok := gl.lineRangeForPolyline(nearestIndex)
	if !ok {
		gl.SetHighlightedLineRange(-1, -1)
		return -1, math.Inf(1)
	}

	gl.SetHighlightedLineRange(start, end)
	return nearestIndex, distance
}

// HighlightNearestPolylineAt converts a world position to Point and highlights
// the nearest polyline.
func (gl *GeometryLayer) HighlightNearestPolylineAt(p fyne.Position) (nearestIndex int, distance float64) {
	return gl.HighlightNearestPolyline(Point{X: float64(p.X), Y: float64(p.Y)})
}

func (gl *GeometryLayer) lineRangeForPolyline(index int) (start int, end int, ok bool) {
	if index < 0 || index >= len(gl.Polylines) {
		return -1, -1, false
	}

	lineOffset := 0
	for i, polyline := range gl.Polylines {
		segmentCount := len(polyline.Points) - 1
		if segmentCount < 1 {
			segmentCount = 0
		}

		if i == index {
			if segmentCount == 0 {
				return -1, -1, false
			}
			return lineOffset, lineOffset + segmentCount - 1, true
		}

		lineOffset += segmentCount
	}

	return -1, -1, false
}

// HighlightNearestPoint finds the nearest point in the GeometryLayer to the given World position p
//   - p - WORLD position in pixels (fyne.Position) to which the nearest point will be highlighted.
func (g *GeometryLayer) HighlightNearestPoint(p fyne.Position) {
	if !g.highlightNearest {
		return
	}
	prevNearest := g.nearestPoint
	prevNearestLine := g.nearestLine
	prevLineStart := g.highlightLineStart
	prevLineEnd := g.highlightLineEnd
	query := Point{X: float64(p.X), Y: float64(p.Y)}

	if len(g.Points) == 0 {
		if g.nearestPoint != -1 {
			g.nearestPoint = -1
		}
	} else {

		var minPointDist float64 = 1.0e+20
		for i, point := range g.Points {
			dist := point.SquareDistance(query)
			if dist < minPointDist {
				minPointDist = dist
				g.nearestPoint = i
			}
		}
	}

	if len(g.Lines) == 0 {
		g.nearestLine = -1
		g.highlightLineStart = -1
		g.highlightLineEnd = -1
	} else {
		var minLineDist float64 = 1.0e+38
		for i, line := range g.Lines {
			dist := line.DistanceTo(query)
			if dist < minLineDist {
				minLineDist = dist
				g.nearestLine = i
			}
		}
		g.highlightLineStart = g.nearestLine
		g.highlightLineEnd = g.nearestLine
	}

	if prevNearest != g.nearestPoint || prevNearestLine != g.nearestLine || prevLineStart != g.highlightLineStart || prevLineEnd != g.highlightLineEnd {
		g.invalidate()
	}
}

// invalidate requests a redraw through both widget and canvas refresh paths.
// This avoids depending on external access to renderer internals.
func (gl *GeometryLayer) invalidate() {
	gl.Refresh()
	canvas.Refresh(gl)
}
