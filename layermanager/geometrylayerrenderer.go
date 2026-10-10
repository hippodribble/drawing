package layermnanager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type GeometryLayerRenderer struct {
	layer   *GeometryLayer
	lines   []*canvas.Line
	points  []*canvas.Circle
	objects []fyne.CanvasObject
}

// Layout positions the lines and points according to the current
// projection. It uses the Sync method to ensure that the number
// of lines and points matches the data in the layer, and then it
// calculates the screen positions of each line and point based on
// their world coordinates and the current projection. The lines
// are positioned using their endpoints, while the points are
// positioned using their center coordinates, adjusted for the
// stroke width.
// IN LAYOUT, WE DO NOT CALL REFRESH - THE RENDERER'S REFRESH CALLS LAYOUT!
func (r *GeometryLayerRenderer) Layout(size fyne.Size) {

	for i, line := range r.layer.Lines {
		p1 := r.layer.Projection.ToScreen(fyne.NewPos(float32(line.X1), float32(line.Y1)))
		p2 := r.layer.Projection.ToScreen(fyne.NewPos(float32(line.X2), float32(line.Y2)))
		r.lines[i].Position1 = fyne.NewPos(p1.X, p1.Y)
		r.lines[i].Position2 = fyne.NewPos(p2.X, p2.Y)
	}

	for i, point := range r.layer.Points {
		p := r.layer.Projection.ToScreen(fyne.NewPos(float32(point.X), float32(point.Y)))
		r.points[i].Move(fyne.NewPos(p.X-r.layer.StrokeWidth/2, p.Y-r.layer.StrokeWidth/2))
		r.points[i].Resize(fyne.NewSize(r.layer.StrokeWidth, r.layer.StrokeWidth))
	}
}

func (r *GeometryLayerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

// Refresh synchronises the canvas objects to match the current state
// of the layer. Then it calls Layout to position the objects according to the current projection.
func (r *GeometryLayerRenderer) Refresh() {
	r.Sync()
	r.Layout(r.layer.Size())
}

func (r *GeometryLayerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *GeometryLayerRenderer) Destroy() {}

// Sync checks if the number of lines and points has changed,
// and if so, it creates new canvas objects for them. It also
// updates the stroke color and width for all lines and points.
// It synchronises the objects with the canvas shapes.
func (r *GeometryLayerRenderer) Sync() {
	nLines := len(r.layer.Lines)
	nPoints := len(r.layer.Points)

	changed := nLines != len(r.lines) || nPoints != len(r.points)

	for len(r.lines) < nLines {
		r.lines = append(r.lines, canvas.NewLine(r.layer.StrokeColor))
	}
	for i := nLines; i < len(r.lines); i++ {
		r.lines[i] = nil
	}
	r.lines = r.lines[:nLines]

	for len(r.points) < nPoints {
		r.points = append(
			r.points,
			canvas.NewCircle(r.layer.StrokeColor),
		)
	}
	for i := nPoints; i < len(r.points); i++ {
		r.points[i] = nil
	}
	r.points = r.points[:nPoints]

	for _, line := range r.lines {
		line.StrokeColor = r.layer.StrokeColor
		line.StrokeWidth = r.layer.StrokeWidth
	}
	for _, point := range r.points {
		point.FillColor = r.layer.StrokeColor
		point.StrokeWidth = r.layer.StrokeWidth / 2
		point.StrokeColor = r.layer.StrokeColor
		point.FillColor = r.layer.fillColor
	}

	if changed {
		r.objects = make(
			[]fyne.CanvasObject,
			0,
			nLines+nPoints,
		)

		// Within this layer: lines underneath point markers.
		for _, line := range r.lines {
			r.objects = append(r.objects, line)
		}
		for _, point := range r.points {
			r.objects = append(r.objects, point)
		}
	}

}

// Sync2 is dumb. It creates all lines and circles and object
// references every time, and does not reuse them.
func (r *GeometryLayerRenderer) Sync2() {

	nLines := len(r.layer.Lines)
	nPoints := len(r.layer.Points)
	r.lines = make([]*canvas.Line, nLines)
	for i := range r.lines {
		r.lines[i] = canvas.NewLine(r.layer.StrokeColor)
		r.lines[i].StrokeWidth = r.layer.StrokeWidth
	}
	r.points = make([]*canvas.Circle, nPoints)
	for i := range r.points {
		r.points[i] = canvas.NewCircle(r.layer.StrokeColor)
		r.points[i].StrokeWidth = r.layer.StrokeWidth / 2
		r.points[i].StrokeColor = r.layer.StrokeColor
		r.points[i].FillColor = r.layer.fillColor
	}
	r.objects = make([]fyne.CanvasObject, nLines+nPoints)
	for i, line := range r.lines {
		r.objects[i] = line
	}
	for i, point := range r.points {
		r.objects[nLines+i] = point
	}

}
