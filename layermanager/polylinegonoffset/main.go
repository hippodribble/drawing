package main

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hippodribble/drawing/layermanager"
	"golang.org/x/image/colornames"
)

var statuschannel chan string
var statuslabel *widget.Label
var manager *layermanager.LayerManager
var centre *fyne.Container

func main() {
	ap := app.New()
	w := ap.NewWindow("Hello")
	w.SetContent(gui())
	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()
}

func gui() fyne.CanvasObject {
	statuslabel = widget.NewLabel("ready")
	statuschannel = make(chan string)
	go status()
	manager = layermanager.NewLayerManager()
	manager.SetStatusChannel(statuschannel)
	top := container.NewHBox(
		widget.NewButton("Polyline", func() {
			fmt.Println(manager.Limits())
			addPolyline(manager)
		}),
	)
	bottom := container.NewStack(statuslabel)

	r := canvas.NewRectangle(color.Transparent)
	r.CornerRadius = 10
	r.SetMinSize(fyne.NewSize(100, 100))
	r.StrokeColor = theme.Color(theme.ColorNameForeground)
	r.StrokeWidth = 1
	centre = container.NewStack(r, manager)
	return container.NewBorder(top, bottom, nil, nil, centre)

}

func status() {
	for update := range statuschannel {
		fyne.Do(func() { statuslabel.SetText(update) })
	}
}
func addPolyline(manager *layermanager.LayerManager) {
	layer := layermanager.NewGeometryLayer()
	layer.SetStrokeColor(colornames.Red)
	layer.SetStrokeWidth(3)

	const (
		radius = 100.0
		cx     = 300.0
		cy     = 300.0
	)

	N:=100
	points := make([]layermanager.Point, 0, N)
	for i := 0; i < N; i++ {
		a := float64(i) * 2 * math.Pi / float64(N)
		points = append(points, layermanager.Point{
			X: cx + radius*math.Cos(a),
			Y: cy + radius*math.Sin(a),
		})
	}
	polygon := layermanager.Polyline{Points: points, ID: "Test Polyline"}.Close()
	polygon=polygon.Reverse()
	layer.AddPolygon(polygon)
	layer.AttachTo(manager)

	expanded := polygon.BufferBy(20)
	shrunken := polygon.BufferBy(-20)
	layer2 := layermanager.NewGeometryLayer()
	layer2.SetStrokeColor(colornames.Blue)
	layer2.SetStrokeWidth(1)
	layer2.AddPolygon(shrunken)
	layer2.AttachTo(manager)

	layer3 := layermanager.NewGeometryLayer()
	layer3.SetStrokeColor(colornames.Green)
	layer3.SetStrokeWidth(1)
	layer3.AddPolygon(expanded)
	layer3.AttachTo(manager)

	rotated:=polygon.Rotate(.1,polygon.Centroid().Translate(50,50))
	layer4 := layermanager.NewGeometryLayer()
	layer4.SetStrokeColor(colornames.Orange)
	layer4.SetStrokeWidth(4)
	layer4.AddPolygon(rotated)
	layer4.AttachTo(manager)

	manager.ShowAll()
}
