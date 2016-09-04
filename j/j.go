package j

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type any interface{}
type array []interface{}

// Node represents json response body
type Node struct {
	any
}

// Array represents json Array
type Array struct {
	array
}

// Property represents property in JSON object
type Property struct {
	Name  string
	Value interface{}
}

// New constructs a Node from data created by JSON serializer
func New(data interface{}) *Node {
	return &Node{data}
}

// Obj constructs a map
func Obj(properties ...*Property) map[string]interface{} {
	var o = make(map[string]interface{})

	for _, p := range properties {
		o[p.Name] = p.Value
	}

	return o
}

// Arr constructs a slice
func Arr(items ...interface{}) []interface{} {
	return items
}

// Prop constructs a json property to be used in Obj function
func Prop(name string, value interface{}) *Property {
	return &Property{
		Name:  name,
		Value: value,
	}
}

// AsArray casts node to an json array
func (n *Node) AsArray() *Array {
	v, ok := n.any.([]interface{})
	if !ok {
		return nil
	}

	return &Array{v}
}

// Raw returns underlying data
func (n *Node) Raw() interface{} {
	return n.any
}

// Path traverses the json and returns a descendent node
func (n *Node) Path(parts ...interface{}) *Node {
	if n == nil {
		return nil
	}

	r := n.any

	for _, part := range parts {
		switch v := part.(type) {
		case int:
			// test for array and do a bounds check
			if array, ok := r.([]interface{}); ok && v >= 0 && v < len(array) {
				r = array[v]
			} else {
				return nil
			}
		case string:
			obj, ok := r.(map[string]interface{})
			if !ok {
				return nil
			}

			tmp, ok := obj[v]
			if !ok {
				return nil
			}

			r = tmp
		default:
			panic("Only string and int supported for indexing")
		}
	}

	return &Node{r}
}

// Find looks for first item that satisfies test
func (a *Array) Find(test func(node *Node) bool) *Node {
	if a == nil {
		return nil
	}

	for _, v := range a.array {
		n := &Node{v}
		if test(n) {
			return n
		}
	}

	return nil
}

// Len returns the length of the array
func (a *Array) Len() int {
	return len(a.array)
}

// Equals checks if value of this node is equal to the paramater
func (n *Node) Equals(val interface{}) bool {
	if v, ok := val.(int); ok {
		return n.any == float64(v)
	}

	return n.any == val
}

// Number attemps to extract a number from the node
func (n *Node) Number() (float64, error) {
	if n == nil {
		return 0, errors.New("Node not found")
	}

	if v, ok := n.any.(float64); ok {
		return v, nil
	}

	return 0, fmt.Errorf("Node is not a number, it's %v", reflect.TypeOf(n.any))
}

func (n *Node) String() (string, error) {
	if n == nil {
		return "", errors.New("Node not found")
	}

	if v, ok := n.any.(string); ok {
		return v, nil
	}

	return "", fmt.Errorf("Node is not a string, it's %v", reflect.TypeOf(n.any))
}

// Dump pretty prints the parsed Json
func (n *Node) Dump() {
	spew.Dump(n.any)
}
