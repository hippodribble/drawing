package drawing

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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

// GeometryLayer is a widget that represents a layer of geometric shapes (points and lines) that can be drawn on a canvas. It holds a collection of points and lines, along with a projection for converting between world and screen coordinates. The layer can be rendered with specified stroke width and color.
type GeometryLayer struct {
	widget.BaseWidget
	Points        []Point
	Lines         []Line
	Projection    *Projection
	statusChannel chan string
	StrokeWidth   float32
	StrokeColor   color.Color
}

// NewGeometryLayer creates a new GeometryLayer with the given projection. It initializes the layer with empty slices for points and lines, and sets the stroke width and color to default values.
func NewGeometryLayer(proj *Projection) *GeometryLayer {
	layer := &GeometryLayer{
		Projection: proj,
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
	gl.Refresh()
}

// SetStatusChannel sets the status channel for the GeometryLayer. The status channel can be used to send status updates or messages related to the layer's operations, such as drawing or refreshing. This is generally used to notify the containing application of relevant changes.
func (gl *GeometryLayer) SetStatusChannel(ch chan string) {
	gl.statusChannel = ch
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
	gl.Refresh()
}
