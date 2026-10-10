# LayerManager package

Package `layermanager` provides a small set of Fyne widgets for drawing simple geometric data as stacked overlay layers.

It is designed for applications that need to display points, lines, and annotations on a shared coordinate system while supporting zooming, fitting the content to view, and layering multiple datasets together.

## Overview

The package contains three core concepts:

- `Point` and `Line` hold world-space geometry
- `GeometryLayer` renders a collection of points and line segments
- `LayerManager` stacks multiple layers and manages a shared projection

## Data model

World-space data are stored as basic primitives:

- `Point` contains `X` and `Y`
- `Line` contains `X1`, `Y1`, `X2`, and `Y2`

A polyline is represented by multiple `Line` values, which keeps the model simple and compatible with Fyne's basic rendering primitives.

## GeometryLayer

`GeometryLayer` is the fundamental drawable unit. Each layer keeps:

- `Points` and `Lines`
- a `Projection` for converting world coordinates to screen coordinates
- `StrokeWidth` and `StrokeColor` for rendering style
- an optional status channel for feedback

The layer will not draw visible points or lines unless a stroke width and stroke color are configured. In practice, set them before or after adding geometry, for example with `SetStrokeWidth(...)` and `SetStrokeColor(...)`.

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

## PolylineLayer

`PolylineLayer` is a convenience layer for working with connected paths rather than a list of independent line segments. It embeds `GeometryLayer` and stores a slice of `Polyline` values, where each `Polyline` is just a sequence of `Point` values.

This makes it easy to model roads, tracks, contours, or other connected shapes without manually splitting each path into many `Line` entries yourself.

When you add a polyline with `AddPolyline`, the layer automatically converts each consecutive pair of points into a segment and appends those segments to the underlying `Lines` slice. This means the polyline still renders through the same geometry pipeline used by `GeometryLayer`.

```go
layer := layermanager.NewPolylineLayer()
polyline := layermanager.NewPolyline([]layermanager.Point{
    {X: 0, Y: 0},
    {X: 5, Y: 2},
    {X: 10, Y: 2},
})

layer.AddPolyline(polyline)
```

`PolylineLayer` also keeps the original `Polylines` slice, so higher-level code can inspect or manipulate the route as a connected path while preserving the rendered line segments for display.

## LayerManager

`LayerManager` combines multiple `GeometryLayer` instances into a single interactive canvas. It manages the shared projection, renders the layers in order, and provides helpers for:

- adding or removing layers
- clearing the stack
- showing the full extent of all content
- zooming around the mouse position
- working in draw mode for interactive sketching

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

