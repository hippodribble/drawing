package drawing

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type AttachableLayer interface {
    AttachTo(*LayerManager)
}

// LayerManager is a widget that manages multiple GeometryLayer instances, allowing for the organization and rendering of geometric shapes on a canvas. It provides functionality for adding, removing, and clearing layers, as well as managing the projection used for coordinate transformations. The LayerManager also handles mouse events for drawing and zooming, and it can communicate status updates through a channel.
type LayerManager struct {
	widget.BaseWidget

	projection    *Projection
	Layers        []*GeometryLayer
	stack         *fyne.Container
	statusChannel chan string
	margin        float32
	startpos      fyne.Position
	mousedown     bool
	drawMode      bool
}

// NewLayerManager creates a new LayerManager with the given projection and an optional list of GeometryLayer instances. It initializes the LayerManager with the specified projection, sets up the stack container for rendering the layers, and prepares it for interaction.
func NewLayerManager(proj *Projection, layers ...*GeometryLayer) *LayerManager {
	v := &LayerManager{
		projection: proj,
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

// SetDrawMode sets the draw mode for the LayerManager. When draw mode is enabled, the LayerManager may allow for interactive drawing on the layers.
func (v *LayerManager) SetDrawMode(draw bool) {
	v.drawMode = draw
}

// RefreshLayers refreshes all GeometryLayer instances managed by the LayerManager, ensuring that they are updated with the current projection and redrawn on the canvas.
func (v *LayerManager) RefreshLayers() {
	for _, layer := range v.Layers {
		layer.Projection = v.projection
		layer.Refresh()
	}
}

// AddLayer adds a new GeometryLayer to the LayerManager, setting its projection to the current projection and updating the stack container to include the new layer.
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

}

func (v *LayerManager) MouseIn(e *desktop.MouseEvent) {}

func (v *LayerManager) MouseMoved(e *desktop.MouseEvent) {
	p := e.Position
	q := v.projection.ToWorld(p)
	if v.statusChannel != nil {
		v.statusChannel <- fmt.Sprintf("World: %9.2f, %9.2f  Screen: %9.2f, %9.2f", q.X, q.Y, p.X, p.Y)
	}
	if v.mousedown && v.drawMode {
		if len(v.Layers) == 0 {
			return
		}
		l := v.Layers[len(v.Layers)-1]
		currentPos := v.projection.ToWorld(p)
		mouseline := NewLine(float64(v.startpos.X), float64(v.startpos.Y), float64(currentPos.X), float64(currentPos.Y))
		l.Lines = []Line{mouseline}
		l.Refresh()
		l.Refresh()
	}
	if !v.mousedown && v.drawMode {
		if len(v.Layers) == 0 {
			return
		}
		// l := v.Layers[len(v.Layers)-1]

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
}

func (v *LayerManager) MouseUp(e *desktop.MouseEvent) {
	v.mousedown = false
	l := v.Layers[len(v.Layers)-1]
	currentPos := v.projection.ToWorld(e.Position)
	mouseline := NewLine(float64(v.startpos.X), float64(v.startpos.Y), float64(currentPos.X), float64(currentPos.Y))
	l.Lines = append(l.Lines, mouseline)
}
