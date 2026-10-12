package layermanager

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeDraw
	ModeMeasure
	ModeIdentify
)

// LayerManager is a widget that manages multiple GeometryLayer instances, allowing for the organization and rendering of geometric shapes on a canvas. It provides functionality for adding, removing, and clearing layers, as well as managing the projection used for coordinate transformations. The LayerManager also handles mouse events for drawing and zooming, and it can communicate status updates through a channel.
type LayerManager struct {
	widget.BaseWidget

	projection    *Projection
	Layers        []*GeometryLayer
	dragLayer     *GeometryLayer
	stack         *fyne.Container
	highlight     bool
	statusChannel chan string
	margin        float32
	startpos      fyne.Position
	mousedown     bool
	mode          Mode
}

// NewLayerManager creates a new LayerManager with the given projection and an optional list of GeometryLayer instances. It initializes the LayerManager with the specified projection, sets up the stack container for rendering the layers, and prepares it for interaction.
func NewLayerManager(layers ...*GeometryLayer) *LayerManager {
	v := &LayerManager{
		projection: NewUnitProjection(),
		Layers:     layers,
		margin:     .1,
	}
	objects := make([]fyne.CanvasObject, len(layers))
	for i, layer := range layers {
		layer.Projection = v.projection
		objects[i] = layer
	}
	v.stack = container.NewStack(objects...)
	v.ExtendBaseWidget(v)
	return v
}

// CreateRenderer creates a new renderer for the LayerManager. The renderer is responsible for drawing the stack of GeometryLayer instances on the canvas, allowing for the visualization of the geometric shapes managed by the LayerManager.
func (v *LayerManager) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.stack)
}

func (v *LayerManager) SetHighlight(highlight bool) {
	v.highlight = highlight
	for _, layer := range v.Layers {
		layer.SetHighlightNearest(highlight)
	}
}

// SetProjection sets the projection for the LayerManager and refreshes all layers to reflect the new projection. The projection is used to convert between world and screen coordinates when rendering the geometric shapes in the layers.
func (v *LayerManager) SetProjection(proj *Projection) {
	v.projection = proj
	v.RefreshLayers()
}

// Projection returns the current projection used by the LayerManager.
func (v *LayerManager) Projection() *Projection {
	return v.projection
}

// SetStatusChannel sets the status channel for the LayerManager and propagates it to all contained GeometryLayer instances. The status channel can be used to send status updates or messages related to the operations of the LayerManager and its layers, such as drawing or refreshing.
func (v *LayerManager) SetStatusChannel(ch chan string) {
	v.statusChannel = ch
	for _, layer := range v.Layers {
		layer.SetStatusChannel(ch)
	}
}

// SetMode sets the interaction mode for the LayerManager, which can be used to control how mouse events are handled. 
// The mode can be set to different values (e.g., ModeNone, ModeDraw, ModeMeasure, ModeIdentify) to enable or disable 
// specific behaviors such as drawing, measuring distances, or identifying geometric shapes. When the mode is set to 
// ModeIdentify, the LayerManager enables highlighting of the nearest points in the layers.
func (v *LayerManager) SetMode(m Mode) {
	v.mode = m
	if m == ModeIdentify {
		v.highlight = true
		for _, layer := range v.Layers {
			layer.SetHighlightNearest(v.highlight)
		}
	} else {
		v.highlight = false
		for _, layer := range v.Layers {
			layer.SetHighlightNearest(v.highlight)
		}
	}
}

// RefreshLayers refreshes all GeometryLayer instances managed by the LayerManager, ensuring that they are updated with the current projection and redrawn on the canvas.
func (v *LayerManager) RefreshLayers() {
	for _, layer := range v.Layers {
		layer.Projection = v.projection
	}
	v.Refresh()
}

// AddLayer adds a new GeometryLayer to the LayerManager, setting its projection to the current projection and updating the stack container to include the new layer.
// deprecated: Use AttachTo method of GeometryLayer instead.
func (v *LayerManager) AddLayer(layer *GeometryLayer) {
	layer.Projection = v.projection
	v.Layers = append(v.Layers, layer)
	v.stack.Objects = append(v.stack.Objects, layer)
	v.Refresh()
}

// RemoveLayer removes the specified GeometryLayer from the LayerManager and updates the stack container accordingly.
func (v *LayerManager) RemoveLayer(layer *GeometryLayer) {
	for i, l := range v.Layers {
		if l == layer {
			v.Layers = append(v.Layers[:i], v.Layers[i+1:]...)
			break
		}
	}
	v.stack.Objects = nil
	for _, layer := range v.Layers {
		v.stack.Objects = append(v.stack.Objects, layer) // ??????
	}
	v.stack.Refresh()
}

// ClearLayers removes all GeometryLayer instances from the LayerManager and updates the stack container accordingly.
func (v *LayerManager) ClearLayers() {
	v.Layers = []*GeometryLayer{}
	v.stack.Objects = []fyne.CanvasObject{}
	v.stack.Refresh()
}

// Limits calculates and returns the minimum and maximum X and Y coordinates of all GeometryLayer instances managed by the LayerManager. It returns four float32 values: minX, minY, maxX, and maxY. If there are no layers, it returns default limits of (0, 0, 0, 0).
// This function is used by the ShowAll() function to fit all layers within the current screen bounds.
func (v *LayerManager) Limits() (minX, minY, maxX, maxY float32) {
	if len(v.Layers) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY, maxX, maxY = v.Layers[0].Limits()
	for _, layer := range v.Layers[1:] {
		lMinX, lMinY, lMaxX, lMaxY := layer.Limits()
		if lMinX < minX {
			minX = lMinX
		}
		if lMinY < minY {
			minY = lMinY
		}
		if lMaxX > maxX {
			maxX = lMaxX
		}
		if lMaxY > maxY {
			maxY = lMaxY
		}
	}
	return
}

// ShowAll adjusts the projection of the LayerManager to ensure that all geometric shapes from all layers are visible within the current screen bounds. It calculates the bounding box of all layers, determines the appropriate scale and center position, and updates the projection accordingly.
func (v *LayerManager) ShowAll() {
	minX, minY, maxX, maxY := v.Limits()

	worldWidth := maxX - minX
	worldHeight := maxY - minY

	screenWidth := v.Size().Width
	screenHeight := v.Size().Height

	scaleX := screenWidth / worldWidth
	scaleY := screenHeight / worldHeight

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	scale *= (1 - v.margin)

	worldCenterX := (minX + maxX) / 2
	worldCenterY := (minY + maxY) / 2

	screenCenterX := screenWidth / 2
	screenCenterY := screenHeight / 2

	mousePos := fyne.NewPos(screenCenterX, screenCenterY)
	worldPos := fyne.NewPos(worldCenterX, worldCenterY)

	proj := NewProjection(mousePos, worldPos, scale)
	v.SetProjection(proj)
	v.Report(fmt.Sprintf("ShowAll: minX=%f, minY=%f, maxX=%f, maxY=%f, scale=%f", minX, minY, maxX, maxY, scale))

}

// Report sends a status update message to the LayerManager's status channel, if it is set.
// This function can be used by the LayerManager or its layers to communicate important
// information or updates to the user interface. Yes, an AI wrote that.
func (v *LayerManager) Report(s string) {
	if v.statusChannel != nil {
		v.statusChannel <- s
	}
}

func (v *LayerManager) MouseIn(e *desktop.MouseEvent) {}

func (v *LayerManager) MouseMoved(e *desktop.MouseEvent) {
	p := e.Position
	q := v.projection.ToWorld(p)
	v.Report(fmt.Sprintf("Position: %9.2f, %9.2f (%9.2f m/pixel)", q.X, q.Y, 1/v.projection.scale))
	if v.mousedown {
		switch v.mode {
		case ModeNone:
			// we must be panning, so adjust the projection based on the mouse movement
			proj := NewProjection(p, v.startpos, v.projection.scale)
			v.SetProjection(proj)
		case ModeMeasure:
			// we must be measuring so draw a line from the start position to the current mouse position
			v.dragLayer.Clear()
			line := NewLine(float64(v.startpos.X), float64(v.startpos.Y), float64(q.X), float64(q.Y))
			v.dragLayer.Lines = append(v.dragLayer.Lines, line)
			v.dragLayer.Refresh()
			v.Report(fmt.Sprintf("Length: %9.2f m", line.Length()))
		}
	}
	if v.mode == ModeIdentify {
		for _, layer := range v.Layers {
			layer.HighlightNearestPoint(q)
			i, _ := layer.HighlightNearestPolylineAt(q)
			if i >= 0 {
				v.Report(layer.Polylines[i].ID)
			}
		}
	}
}

func (v *LayerManager) Scrolled(e *fyne.ScrollEvent) {
	d := e.Scrolled
	if d.DY > 0 {
		v.Zoom(1.05, e.Position)
	} else if d.DY < 0 {
		v.Zoom(0.95, e.Position)
	}
}

func (v *LayerManager) Zoom(factor float32, mouse fyne.Position) {
	proj := v.projection
	world := proj.ToWorld(mouse)
	newScale := proj.scale * factor
	newProj := NewProjection(mouse, world, newScale)
	v.SetProjection(newProj)
}

func (v *LayerManager) MouseOut() {}

func (v *LayerManager) MouseDown(e *desktop.MouseEvent) {
	v.startpos = v.projection.ToWorld(e.Position)
	v.mousedown = true
	if v.mode == ModeMeasure { // create a drag layer if necessary
		if v.dragLayer == nil {
			v.dragLayer = NewGeometryLayer()
			v.dragLayer.StrokeColor = colornames.Red
			v.dragLayer.StrokeWidth = 2.0
			v.AddLayer(v.dragLayer)
		}
		v.dragLayer.Clear()
	}
}

func (v *LayerManager) MouseUp(e *desktop.MouseEvent) {
	v.mousedown = false
	v.RemoveLayer(v.dragLayer)
	v.dragLayer = nil
	if e.Button == desktop.MouseButtonSecondary {
		v.ShowAll()
	}

}
