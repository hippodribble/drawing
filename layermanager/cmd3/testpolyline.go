package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	drawing "github.com.hippodribble/drawing/layermanager"
)

var winMain fyne.Window

func main() {
	ap := app.NewWithID("canvas.test")
	winMain = ap.NewWindow("Canvas Test")
	winMain.Resize(fyne.NewSize(1500, 1000))
	winMain.SetContent(gui())
	winMain.ShowAndRun()
}

func gui() fyne.CanvasObject {
	// Create a new polyline layer
	polylineLayer := &drawing.PolylineLayer{}

	// Add a sample polyline
	polyline := drawing.Polyline{
		Points: []drawing.Point{
			{X: 0, Y: 0},
			{X: 100, Y: 100},
			{X: 200, Y: 50},
		},
	}
	polylineLayer.AddPolyline(polyline)
	polylineLayer.SetStrokeWidth(2)
	polylineLayer.SetStrokeColor(theme.Color(theme.ColorNameForeground))

	layermanager := drawing.NewLayerManager(drawing.NewUnitProjection())

	top:=widget.NewButton("Draw",func() {
	polylineLayer.AttachTo(layermanager)
	layermanager.RefreshLayers()
	layermanager.ShowAll()

	})


	// Return the polyline layer as the content of the window

	return container.NewBorder(top, nil, nil, nil, container.NewStack(layermanager))
}
