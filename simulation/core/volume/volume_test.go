package volume

import (
	"testing"
)

var v Volume

func init() {
	v = NewVolume(10, 10, 0.)
}

func TestVolumeClone(t *testing.T) {
	cv := v.Clone()
	cv.Set(0, 0, 99.)

	before := v.Get(0, 0)
	after := cv.Get(0, 0)
	if !(before == 0. && after == 99.) {
		t.Errorf("before: %f after: %f", before, after)
	}
}

func TestSliceCopy(t *testing.T) {
	v1 := 1.
	v2 := 999.

	s := [][]float64{{v1, v1, v1}, {v1, v1, v1}}
	var cs [][]float64
	for i := 0; i < len(s); i++ {
		line := make([]float64, len(s[0]))
		copy(line, s[i])
		cs = append(cs, line)
	}
	cs[0][0] = v2

	if !(s[0][0] == v1 && cs[0][0] == v2) {
		t.Error("slice: ", s, "clone: ", cs)
	}
}
