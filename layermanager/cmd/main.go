package main

import (
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	layermanager "github.com/hippodribble/drawing/layermanager"

	"golang.org/x/image/colornames"
)

var viewport *layermanager.LayerManager
var dragLayer *layermanager.GeometryLayer
var top *fyne.Container
var centre *fyne.Container
var status = widget.NewLabel("Status: Ready")
var bottom = container.NewVBox(status)
var statusChannel = make(chan string)

func main() {
	ap := app.NewWithID("canvas.test")
	w := ap.NewWindow("Canvas Test")
	w.Resize(fyne.NewSize(1500, 1000))
	w.SetContent(gui())
	w.ShowAndRun()
}

func gui() fyne.CanvasObject {
	go monitorStatusUpdates()
	viewport = layermanager.NewLayerManager()
	viewport.SetStatusChannel(statusChannel)
	dragLayer = layermanager.NewGeometryLayer()
	dragLayer.StrokeColor = theme.Color(theme.ColorNameForeground)
	dragLayer.StrokeWidth = 2.0
	viewport.AddLayer(dragLayer)
	centre = container.NewStack(viewport)
	top = container.NewHBox(
		widget.NewButton("Point Layer", addDummyPointLayer),
		widget.NewButton("Line Layer", addDummyLineLayer),
		widget.NewButton("Point+Line Layer", addDummyPointLineLayer),
		widget.NewButton("Clear Layers", viewport.ClearLayers),
		widget.NewButton("Show All", viewport.ShowAll),

	)
	b := container.NewBorder(
		top,
		bottom,
		nil,
		nil,
		viewport,
	)
	return b
}

func monitorStatusUpdates() {
	for update := range statusChannel {
		fyne.Do(func() { status.SetText(update) })
	}
}

func randColor() color.Color {
	colors := []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
		colornames.Orange,
		colornames.Purple,
		colornames.Cyan,
		colornames.Magenta,
		colornames.Yellow,
	}
	return colors[rand.Intn(len(colors))]
}

func randColor2() color.Color {
	r := rand.Uint32() >> 24
	g := rand.Uint32() >> 24
	b := rand.Uint32() >> 24

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func addDummyPointLayer() {
	for range 100 {
		layer := makeDummyPointLayer(100)
		layer.AttachTo(viewport)
	}
	viewport.ShowAll()
	viewport.RefreshLayers()
}

func addDummyLineLayer() {
	layer := makeDummyLineLayer(rand.Intn(10) + 5)
	layer.AttachTo(viewport)
	viewport.ShowAll()
	viewport.RefreshLayers()
}

func addDummyPointLineLayer() {
	layer := makeDummyPointLineLayer(rand.Intn(10) + 5)
	layer.AttachTo(viewport)
	viewport.ShowAll()
	viewport.RefreshLayers()
}

func makeDummyPointLayer(numPoints int) *layermanager.GeometryLayer {
	layer := layermanager.NewGeometryLayer()
	x0 := rand.NormFloat64() * 20
	y0 := rand.NormFloat64() * 10
	for i := 0; i < numPoints; i++ {
		x := rand.NormFloat64() + x0
		y := rand.NormFloat64() + y0
		point := layermanager.NewPoint(x, y)
		layer.Points = append(layer.Points, point)
	}
	layer.StrokeWidth = 2
	layer.StrokeColor = randColor2()
	return layer
}

func makeDummyLineLayer(numLines int) *layermanager.GeometryLayer {
	layer := layermanager.NewGeometryLayer()
	for i := 0; i < numLines; i++ {
		x1 := rand.NormFloat64()
		y1 := rand.NormFloat64()
		x2 := rand.NormFloat64()
		y2 := rand.NormFloat64()
		line := layermanager.NewLine(x1, y1, x2, y2)
		layer.Lines = append(layer.Lines, line)
	}
	layer.StrokeWidth = 1
	layer.StrokeColor = randColor()
	return layer
}

func makeDummyPointLineLayer(numPointsLines int) *layermanager.GeometryLayer {
	layer := layermanager.NewGeometryLayer()
	for i := 0; i < numPointsLines; i++ {
		x := rand.NormFloat64()
		y := rand.NormFloat64()
		point := layermanager.NewPoint(x, y)
		layer.Points = append(layer.Points, point)

		x1 := rand.NormFloat64()
		y1 := rand.NormFloat64()
		x2 := rand.NormFloat64()
		y2 := rand.NormFloat64()
		line := layermanager.NewLine(x1, y1, x2, y2)
		layer.Lines = append(layer.Lines, line)
	}
	for i := 0; i < numPointsLines; i++ {
		x := rand.NormFloat64()
		y := rand.NormFloat64()
		point := layermanager.NewPoint(x, y)
		layer.Points = append(layer.Points, point)
	}
	layer.StrokeWidth = 1
	layer.StrokeColor = randColor()
	return layer
}
