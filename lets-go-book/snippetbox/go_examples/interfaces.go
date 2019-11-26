package main

import (
	"fmt"
)

type Dog struct {
	Feature string
	Animal
}

func (d Dog) Special() {
	fmt.Println(d.Feature, "The dog barks")
}

type Cat struct {
	Feature string
	Animal
}

func (c Cat) Special() {
	fmt.Println(c.Feature, "The cat purrs")
}

type Animal struct {
	Name  string
	Breed string
}

func (c Animal) Walk() {
	fmt.Println(c.Name, "walks across the room")
}
func (c Animal) Sit() {
	fmt.Println(c.Name, "sits down")
}

// This interface represents any type
// that has both Walk and Sit methods.
type FourLegged interface {
	Walk()
	Sit()
}

// We can replace DemoDog and DemoCat
// with this single function.
func Demo(animal FourLegged) {
	//x, ok := animal.(*Dog)
	if animal == nil { // || ok && x == nil {
		fmt.Println("its nil")
		return
	}

	animal.Walk()
	animal.Sit()

	switch at := animal.(type) {
	case Dog:
		at.Special()
	case Cat:
		at.Special()
	default:

	}

	//at = animal.(Dog) // If its not a Dog, will panic
	//at.Special()

	dog, _ := animal.(Dog)
	dog.Special() // Will return ' '

	dog2, ok := animal.(Dog)
	if ok {
		dog2.Special()
	}

}

func main() {
	//dog := Dog{"Fido", "Terrier"}
	//cat := Cat{"Fluffy", "Siamese"}
	monkey := Animal{"Malumba", "Red Ass Monkey"}
	//Demo(dog)
	// The above call (again) outputs:
	// Fido walks across the room
	// Fido sits down
	//Demo(cat)
	dog := Dog{
		Feature: "dog",
		Animal:  Animal{Name: "Fido", Breed: "Terrier"},
	}
	// The above call (again) outputs:
	// Fluffy walks across the room
	// Fluffy sits down
	Demo(dog)
	Demo(monkey)
	Demo(nil)
	//Demo((*Dog)(nil))
}
