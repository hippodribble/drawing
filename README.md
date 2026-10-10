# Drawing

A small Go library for layered geometric drawing in Fyne applications.

The project provides a lightweight canvas toolkit for visualising point and line data in separate overlay layers. It is designed for applications that need a shared world-to-screen projection, interactive zooming, and composable drawing layers.

## Features

- `GeometryLayer` for simple point and line data
- `LayerManager` for stacking multiple layers in one widget
- shared `Projection` for world/screen coordinate transforms
- mouse wheel zoom and fit-to-view helpers
- status reporting via channels for drawing and interaction updates

## Package layout

- `layermanager` contains the Fyne widget code for layered rendering
- `cmd`, `cmd2`, and `cmd3` contain small example programs

## Installation

```bash
go get github.com/hippodribble/drawing/layermanager
```

## Quick example

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "github.com/hippodribble/drawing/layermanager"
)

func main() {
    a := app.New()
    w := a.NewWindow("Drawing")

    manager := layermanager.NewLayerManager()
    layer := layermanager.NewGeometryLayer()

    layer.Points = []layermanager.Point{
        {X: 0, Y: 0},
        {X: 10, Y: 5},
    }
    layer.Lines = []layermanager.Line{
        {X1: 0, Y1: 0, X2: 10, Y2: 5},
    }

    layer.AttachTo(manager)
    manager.ShowAll()
    w.SetContent(manager)
    w.ShowAndRun()
}
```

## Notes

This project is intentionally simple: geometry is stored as basic point and segment primitives, and the rendering stack focuses on clarity and composability instead of a full CAD-style model.
