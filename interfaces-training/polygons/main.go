package main

import (
	"fmt"
	"math"
)

type WeirdRectangle struct {
	x1, y1, x2, y2 float64
}

type Circle struct {
	radius, diagonal float64
}

type Rectangle struct {
	length, width float64
}

func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

func circleArea(radius, diagonal float64) float64 {
	return math.Pi * radius * radius
}

func rectangleArea(length, width float64) float64 {
	return length * width
}

func weirdRectangleArea(x1, y1, x2, y2 float64) float64 {
	l := distance(x1, y1, x1, y2)
	w := distance(x1, y1, x2, y1)
	return l * w
}

// METHODS
func (c *Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (r *Rectangle) diagonal() float64 {
	return math.Sqrt(r.width*r.width + r.length*r.length)
}

func (r *Rectangle) area() float64 {
	return r.width * r.length
}

// Because both of our polygons have an area method, can be implemented by an
// interface....
type Shape interface {
	area() float64
	// With a Shape instance, we wouldn't be able to access the struct fields
	// for the polygon. Only the area.
}

func totalArea(shapes ...Shape) float64 {
	var area float64
	for _, s := range shapes {
		area += s.area()
	}
	return area
}

type MultiShape struct {
	shapes []Shape
}

type GetArea struct {
	Shape
}

func (m *MultiShape) area() float64 {
	var area float64
	for _, s := range m.shapes {
		area += s.area()
	}
	return area
}

func (a *GetArea) area() float64 {
	return a.Shape.area()
}

func main() {
	var c Circle
	var r Rectangle

	fmt.Println(weirdRectangleArea(2, 2, 4, 4))

	wr := WeirdRectangle{2, 2, 4, 4}
	fmt.Println(weirdRectangleArea(wr.x1, wr.y1, wr.x2, wr.y2))

	c = Circle{5, 0}
	fmt.Println(circleArea(c.radius, c.diagonal))

	c = Circle{5, 0}
	fmt.Println(c.area())
	r = Rectangle{3, 5}
	fmt.Println(r.area())
	fmt.Println(totalArea(&c, &r))

	// In Go, interfaces define functionality, rather than data, so interfaces
	// can also be used as fields...
	//
	// Here, Shape interface is being used as a field.
	multiShape := MultiShape{
		shapes: []Shape{
			&Circle{5, 0},
			&Rectangle{3, 5},
		},
	}
	fmt.Println(multiShape.area())

	polygon1Area := GetArea{
		&Circle{5, 0},
	}
	fmt.Println(polygon1Area.area())
	polygon2Area := GetArea{
		&Rectangle{5, 2},
	}
	fmt.Println(polygon2Area.area())
}
