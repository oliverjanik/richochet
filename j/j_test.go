package j

import (
	"encoding/json"
	"testing"
)

func TestObj(t *testing.T) {
	o := Obj(
		Prop("name", "ricochet"),
		Prop("version", 1),
		Prop("isAwesome", true),
	)

	b, err := json.Marshal(o)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
}

func TestArray(t *testing.T) {
	a := Array(
		true,
		25,
		8.5,
		"text",
		nil,
		Obj(
			Prop("nested", false),
		),
	)

	b, err := json.Marshal(a)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
}
