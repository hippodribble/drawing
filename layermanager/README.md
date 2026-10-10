# Drawing Package

Package _drawing_ provides:
* a simple framework to display points (as circles) and line segments in layer widgets of type `GeometryLayer`
* a `LayerManager` to stack `GeometryLayer` instances in a single widget

## Data
 
World data are either points or line segments, expressed as `Point` or `Line` types.

`Point`s are very basic, with only `X` and `Y` values. This allows adapters to be written to/from user types

`Line`s are represented by `X1`, `Y1`, `X2`, `Y2` values They are necessarily only segments between two points, as that is all Fyne supports. Polylines (open or closed) are represented by multiple `Line`s.

## Geometry Layers ##

`GeometryLayer` is the basic drawable unit in the package. It is a Fyne widget that stores a set of `Point` and `Line` values together with the information needed to render them on screen.

Each layer contains:

* `Points` and `Lines`, which hold the world-space geometry for that layer
* a `Projection`, which converts between world coordinates and screen coordinates during rendering
* `StrokeWidth` and `StrokeColor`, which control the appearance of both points and line segments
* an optional status channel used to report layer-related status back to the application

Layers are intended to be independent overlays. A viewport or layer manager can stack multiple `GeometryLayer` instances so that each dataset, annotation set, or interactive drawing pass remains separate.

The main operations provided by `GeometryLayer` are:

* `NewGeometryLayer`, to construct a new empty layer for a projection
* `SetProjection`, to update the coordinate transform and refresh the display
* `Limits`, to calculate the layer bounds from all points and segment endpoints
* `Clear`, to remove all stored geometry and redraw the layer empty

In practice, a layer is usually populated by appending `Point` or `Line` values and then letting the containing layer manager refresh and render it.

## Layer Manager ##

`LayerManager` is the widget that combines multiple `GeometryLayer` instances into a single interactive canvas. It owns the shared `Projection`, renders the layers as a stack, and handles user input such as mouse movement, button presses, and scroll-wheel zoom.

Its main responsibilities are:

* maintaining an ordered list of layers in `Layers`
* applying a shared projection to every layer before refresh
* providing layer lifecycle operations such as `AddLayer`, `RemoveLayer`, and `ClearLayers`
* reporting cursor position through a status channel in both world and screen coordinates
* fitting all visible geometry into view with `ShowAll`
* supporting interactive drawing when `SetDrawMode(true)` has been enabled

The layer order is significant. Layers are rendered in stack order, and the topmost layer is also the one used by the current interactive drawing logic.

## Projection ##

`Projection` converts between world coordinates and screen coordinates. Rather than storing a fixed origin, it is defined by three values:

* a screen-space reference point (`mouse`)
* the world-space point that corresponds to that screen position (`world`)
* a scale factor (`scale`)

This is enough to perform bidirectional conversion:

* `ToScreen` maps a world point into widget coordinates
* `ToWorld` maps a widget position back into world space

The package also provides a few helpers around this model:

* `NewProjection` builds a projection from explicit reference values
* `NewUnitProjection` creates a safe default projection
* `ZoomIn` and `ZoomOut` create derived projections that preserve the world position under the mouse while changing scale

This design keeps zooming intuitive, because the point under the cursor remains stable while the scale changes.

