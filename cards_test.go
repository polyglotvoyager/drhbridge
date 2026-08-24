package main

import (
    "testing"
)

func TestCardBeats(t *testing.T) {
    a := Card{"K", "S"}
    b := Card{"9", "S"}

    if !CardBeats(a, b, "H", "C") {
        t.Errorf("%q should beat %q", a, b)
    }
}
