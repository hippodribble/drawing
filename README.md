# Drawing Package for fyne.io/fyne/v2

## LayerManager (Fyne widget)

This package is used for layered drawing of points and lines. It includes:

### GeometryLayer

Takes point and line segment data in a simple format, and combines them into a drawing layer

### LayerManager

The layer manager is used to compose multiple layers for display. It handles world / screen units projection, mouse scroll zoom, etc.

### Projection

Represents the screen-world-screen transform as a datum and scale
