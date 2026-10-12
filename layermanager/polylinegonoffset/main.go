package main

import (
	"image/color"
	"math"
	"math/rand/v2"

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
		widget.NewButton("Polygon", func() {
			addPolygon(manager)
		}),
		widget.NewButton("Polyline", func() {
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
func addPolygon(manager *layermanager.LayerManager) {

	manager.ClearLayers()

	layer := layermanager.NewGeometryLayer()
	layer.SetStrokeColor(colornames.Red)
	layer.SetStrokeWidth(3)

	const (
		radius = 100.0
		cx     = 300.0
		cy     = 300.0
	)

	N := 51
	points := make([]layermanager.Point, 0, N)
	for i := 0; i < N; i++ {
		a := float64(i) * 2 * math.Pi / float64(N)
		r := radius + rand.NormFloat64()*2
		points = append(points, layermanager.Point{
			X: cx + r*math.Cos(a),
			Y: cy + r*math.Sin(a),
		})
	}
	polygon := layermanager.Polyline{Points: points, ID: "Test Polyline"}.Close()
	polygon = polygon.Reverse()
	layer.AddPolygon(polygon)
	layer.AttachTo(manager)


	dp:=polygon.DouglasPeucker(10)
	layer5 := layermanager.NewGeometryLayer()
	layer5.SetStrokeColor(colornames.Purple)
	layer5.SetStrokeWidth(2)
	layer5.AddPolygon(dp)
	layer5.AttachTo(manager)

	manager.ShowAll()
}

func addPolyline(manager *layermanager.LayerManager) {

	manager.ClearLayers()

	points := []layermanager.Point{}
	for i := range 1000 {
		points = append(points, layermanager.Point{
			X: float64(i) * 10,
			Y: float64(i)*5 + rand.NormFloat64()*100,
		})
	}

	polyline := layermanager.Polyline{Points: points, ID: "Test Polyline"}
	layer := layermanager.NewGeometryLayer()
	layer.SetStrokeColor(colornames.Red)
	layer.SetStrokeWidth(1)
	layer.AddPolyline(polyline)
	layer.AttachTo(manager)

	x := 100.0

	polyline2 := polyline.DouglasPeucker(x * 6)
	layer2 := layermanager.NewGeometryLayer()
	layer2.SetStrokeColor(colornames.Black)
	layer2.SetStrokeWidth(2)
	layer2.AddPolyline(polyline2)
	layer2.AttachTo(manager)
	manager.ShowAll()

}
