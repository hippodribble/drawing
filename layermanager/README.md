# LayerManager package

Package `layermanager` provides a small set of Fyne widgets for drawing simple geometric data as stacked overlay layers.

It is designed for applications that need to display points, lines, and annotations on a shared coordinate system while supporting zooming, fitting the content to view, and layering multiple datasets together.

## Overview

The package uses these core concepts:

- `Point`, `Line`, `Polyline` and `Polygon` hold world-space geometry
- `GeometryLayer` renders a collection of points and line segments, and deconstructs `Polyline` and `Polygon` types to their component `Line` segments
- `LayerManager` stacks multiple layers and manages a shared projection
- shape convenience functions are available on `Point`, `Line`, `Polyline`, and `Polygon`

## Data model

World-space data are stored as basic primitives:

- `Point` contains `X` and `Y`
- `Line` contains `X1`, `Y1`, `X2`, and `Y2`
- `Polyline` contains an ordered list of `Point` values
- `Polygon` contains a ordered, closed list of `Point` values

Exported geometry functions in `geometry.go`:

- `Point`
    - `NewPoint(x1, y1 float64) Point`
    - `(Point) Translate(dx, dy float64) Point`
    - `(Point) Rotate(angle float64, centre Point) Point`
    - `(Point) SquareDistance(other Point) float64`
- `Line`
    - `NewLine(x1, y1, x2, y2 float64) Line`
    - `(Line) Translate(dx, dy float64) Line`
    - `(Line) Rotate(angle float64, centre Point) Line`
    - `(Line) Length() float64`
    - `(Line) DistanceTo(p Point) float64`
- `Polyline`
    - `NewPolyline(points []Point) Polyline`
    - `(Polyline) Translate(dx, dy float64) Polyline`
    - `(Polyline) Rotate(angle float64, centre Point) Polyline`
    - `(Polyline) Length() float64`
    - `(Polyline) DouglasPuecker(accuracy float64) Polyline`
    - `(Polyline) Close() Polygon`
    - `(Polyline) Offset(offset float64) Polyline`
- `Polygon`
    - `NewPolygon(points []Point) Polygon`
    - `(Polygon) Close() Polygon`
    - `(Polygon) Translate(dx, dy float64) Polygon`
    - `(Polygon) Rotate(angle float64, centre Point) Polygon`
    - `(Polygon) Area() float64`
    - `(Polygon) Perimeter() float64`
    - `(Polygon) Centroid() Point`
    - `(Polygon) IsClockwise() bool`
    - `(Polygon) Reverse() Polygon`
    - `(Polygon) BufferBy(r float64) Polygon`
    - `(Polygon) DouglasPeucker(tolerance float64) Polygon`

## GeometryLayer

`GeometryLayer` is the fundamental drawable unit. Each layer keeps:

- `Points`, `Lines` and `Polygons`
- a `Projection` for converting world coordinates to screen coordinates
- `StrokeWidth`, `FillColor` and `StrokeColor` for rendering style
- an optional status channel for feedback

The layer will not draw visible points or lines unless a stroke width and stroke color are configured. In practice, set them before or after adding geometry, for example with `SetStrokeWidth(...)` and `SetStrokeColor(...)`. There are defaults.

Layers are intended to be independent overlays that can be added, removed, cleared, or reused as needed.

Typical usage looks like this:

```go
import "image/color"

layer := layermanager.NewGeometryLayer()
layer.SetStrokeWidth(2)
layer.SetStrokeColor(color.Black)
layer.Points = []layermanager.Point{{X: 0, Y: 0}, {X: 10, Y: 5}}
layer.Lines = []layermanager.Line{{X1: 0, Y1: 0, X2: 10, Y2: 5}}
```

## LayerManager

`LayerManager` combines multiple `GeometryLayer` instances into a single interactive canvas. It manages the shared projection, renders the layers in order, and provides helpers for:

- adding or removing layers to a stack (last layer displays on top of earlier layers)
- clearing the stack
- showing the full extent of all content
- zooming around the mouse position (a.k.a. mouse-centred zoom)
- panning the drawing using left mouse button

Interaction modes can be changed with `SetMode(...)`:

- `ModeNone`: default navigation mode. Drag with left mouse button to pan, scroll to zoom, and right-click to run `ShowAll()`.
- `ModeDraw`: reserved mode for drawing workflows. It is currently defined but has no custom handling in mouse event logic.
- `ModeMeasure`: drag with left mouse button to draw a temporary measurement line and report its length.
- `ModeIdentify`: enables nearest-geometry highlighting and reports IDs of nearby polylines while moving the mouse.

```go
manager := layermanager.NewLayerManager()
layer.AttachTo(manager)
manager.ShowAll()
```



## Projection

`Projection` converts between world and screen coordinates. It stores:

- a screen-space reference point (`mouse`)
- the corresponding world-space point (`world`)
- a scale factor (`scale`)

This is enough to support `ToScreen`, `ToWorld`, and zoom operations without losing the point under the cursor.

## Example workflow

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "github.com/hippodribble/drawing/layermanager"
)

func main() {
    a := app.New()
    w := a.NewWindow("Layered drawing")

    manager := layermanager.NewLayerManager()
    layer := layermanager.NewGeometryLayer()

    layer.Points = []layermanager.Point{{X: 0, Y: 0}, {X: 10, Y: 10}}
    layer.Lines = []layermanager.Line{{X1: 0, Y1: 0, X2: 10, Y2: 10}}
    layer.AttachTo(manager)

    manager.ShowAll()
    w.SetContent(manager)
    w.ShowAndRun()
}
```

## Notes

This package intentionally keeps its geometry primitives minimal. It is a good fit for simple overlays, annotation layers, plotting views, and lightweight interactive drawing tools built on Fyne.

The shape functions are meant to keep common geometric work directly on the model objects, so you can translate, rotate, offset, and measure `Point`, `Line`, `Polyline`, and `Polygon` values without extra helper code.

