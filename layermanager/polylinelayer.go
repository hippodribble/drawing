package layermanager

// PolylineLayer is a widget that represents a layer of polylines that can be drawn on a canvas. 
// It extends the GeometryLayer to hold a collection of polylines, which are sequences of connected 
// points. The layer can be rendered with specified stroke width and color, and it provides methods 
// for adding and clearing polylines, as well as calculating the limits of the polylines in the layer.
type PolylineLayer struct {
	*GeometryLayer
	Polylines []Polyline
}

// NewPolylineLayer creates a new PolylineLayer with an embedded GeometryLayer. 
// It initializes the layer with an empty slice for polylines and sets up the 
// necessary properties for rendering and interaction.
// The GeometryLayer must be a reference type, so we use a pointer to it. 
func NewPolylineLayer() *PolylineLayer {
	layer := &PolylineLayer{
		GeometryLayer: NewGeometryLayer(),
	}
	// layer.ExtendBaseWidget(layer)
	return layer
}

// AddPolyline adds a new polyline to the PolylineLayer. It appends the given polyline to the Polylines 
// slice and also converts the points of the polyline into individual lines, which are added to the Lines 
// slice of the GeometryLayer. After adding the polyline, it refreshes the layer to update the rendering.
func(pl *PolylineLayer) AddPolyline(polyline Polyline) {
	pl.Polylines = append(pl.Polylines, polyline)
	for i:=0;i<len(polyline.Points)-1;i++ {
		line := Line{
			X1: polyline.Points[i].X,
			Y1: polyline.Points[i].Y,
			X2: polyline.Points[i+1].X,
			Y2: polyline.Points[i+1].Y,
		}
		pl.Lines = append(pl.Lines, line)
	}
	pl.Refresh()
}

// ClearPolylines removes all polylines from the PolylineLayer. It clears the Polylines slice and also clears the Lines slice of the GeometryLayer, effectively removing all geometric representations of the polylines. After clearing, it refreshes the layer to update the rendering.
func(pl *PolylineLayer) ClearPolylines() {
	pl.Polylines = pl.Polylines[:0]
	pl.Lines = pl.Lines[:0]
	pl.Refresh()
}

// Limits calculates and returns the minimum and maximum X and Y coordinates of all polylines in the PolylineLayer.
// It returns four float32 values: minX, minY, maxX, and maxY. If the layer has no polylines, it returns default limits of (0, 0, 1, 1). 
// This method is useful for determining the bounding box of the polylines in the layer.
func (pl *PolylineLayer) Limits() (minX, minY, maxX, maxY float32) {
	if len(pl.Polylines) == 0 {
		return 0, 0, 1, 1
	}
	minX, maxX, minY, maxY = 1.0e+20, -1.0e+20, 1.0e+20, -1.0e+20
	for _, polyline := range pl.Polylines {
		for _, point := range polyline.Points {
			if float32(point.X) < minX {
				minX = float32(point.X)
			}
			if float32(point.X) > maxX {
				maxX = float32(point.X)
			}
			if float32(point.Y) < minY {
				minY = float32(point.Y)
			}
			if float32(point.Y) > maxY {
				maxY = float32(point.Y)
			}
		}
	}
	return
}
// AttachTo attaches the PolylineLayer to a LayerManager. It adds the PolylineLayer's GeometryLayer to the LayerManager, sets the projection of the PolylineLayer to match that of the LayerManager, and refreshes the LayerManager to ensure that the new layer is properly displayed. This method allows for easy integration of the PolylineLayer into a larger layer management system.
func (pl *PolylineLayer) AttachTo(m *LayerManager) {
	m.AddLayer(pl.GeometryLayer)
	pl.Projection=m.Projection()
	m.Refresh()
}