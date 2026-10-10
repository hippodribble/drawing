package layermanager

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// Projection is a simple struct that holds the parameters for projecting world coordinates to screen coordinates and vice versa. It contains datum (screen/world relationship) and scale information.
// It contains the following fields:
//   - mouse: The position of the mouse on the screen, which serves as a reference point for the projection.
//   - world: The position in the world that corresponds to the mouse position, which is used to determine how to convert between world and screen coordinates.
//   - scale: The scale factor for converting world coordinates to screen coordinates. A larger scale factor means that world coordinates will be scaled up more when converting to screen coordinates, while a smaller scale factor means that they will be scaled down more.
//
// The Projection struct provides methods for converting between world and screen coordinates, as well as for zooming in and out of the projection.
//   - No origin is required - the mouse/world/scale relationship is sufficient to define the projection. The mouse position serves as a reference point for the projection, and the world position indicates the corresponding point in the world coordinates. The scale factor determines how much the world coordinates are scaled when converting to screen coordinates.
//   - when zooming in or out, neither the mouse nor world position changes, only the scale.
type Projection struct {
	mouse fyne.Position // The position of the mouse on the screen
	world fyne.Position // The position in the world that corresponds to the mouse position
	scale float32       // The scale factor for converting world coordinates to screen coordinates
}

// NewProjection creates a new Projection with the given mouse position, world position, and scale factor.
func NewProjection(mouse, world fyne.Position, scale float32) *Projection {
	return &Projection{
		mouse: mouse,
		world: world,
		scale: scale,
	}
}

// NewUnitProjection creates a new Projection with default values - it is just a placeholder and is unlikely to be useful apart from preventing divide-by-zero errors.
func NewUnitProjection() *Projection {
	return &Projection{
		mouse: fyne.NewPos(100, 100),
		world: fyne.NewPos(1000, 1000),
		scale: 1.0,
	}
}

// String returns a string representation of the Projection, including the mouse position, world position, and scale factor.
func (proj *Projection) String() string {
	return fmt.Sprintf("Projection: mouse=(%f,%f), world=(%f,%f), scale=%f", proj.mouse.X, proj.mouse.Y, proj.world.X, proj.world.Y, proj.scale)
}

// ToScreen converts world coordinates to screen coordinates using the current projection parameters.
func (proj *Projection) ToScreen(newWorld fyne.Position) fyne.Position {
	qx := proj.mouse.X + (newWorld.X-proj.world.X)*proj.scale
	qy := proj.mouse.Y - (newWorld.Y-proj.world.Y)*proj.scale
	return fyne.NewPos(qx, qy)
}

// ToWorld converts screen coordinates to world coordinates using the current projection parameters.
func (proj *Projection) ToWorld(screen fyne.Position) fyne.Position {
	sx := screen.X
	sy := screen.Y
	px := proj.world.X + (sx-proj.mouse.X)/proj.scale
	py := proj.world.Y - (sy-proj.mouse.Y)/proj.scale
	return fyne.NewPos(px, py)
}

// ZoomOut creates a new Projection that is zoomed out by the given factor, using the current mouse position as the reference point. The world position is updated to reflect the world position at the mouse location, and the scale factor is multiplied by the given factor.
func (proj *Projection) ZoomOut(m fyne.Position, factor float32) *Projection {
	w := proj.ToWorld(m)
	scale := proj.scale * factor
	return NewProjection(m, w, scale)
}

// ZoomIn creates a new Projection that is zoomed in by the given factor, using the current mouse position as the reference point. The world position is updated to reflect the world position at the mouse location, and the scale factor is divided by the given factor.
func (proj *Projection) ZoomIn(m fyne.Position, factor float32) *Projection {
	w := proj.ToWorld(m)
	scale := proj.scale / factor
	return NewProjection(m, w, scale)
}
