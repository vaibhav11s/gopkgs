package vector2d

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testNewVector(t *testing.T, new func(a interface{}, b interface{}) (Vector2D, error)) {
	tests := []struct {
		a, b interface{}
		v    Vector2D
		err  bool
	}{
		{0, 0, Vector2D{X: 0, Y: 0}, false},
		{1, 2, Vector2D{X: 1, Y: 2}, false},
		{1.2, 2.0, Vector2D{X: 1.2, Y: 2}, false},
		{1.2, "2.1", Vector2D{}, true},
		{"1.2", "2.1", Vector2D{}, true},
		{int32(1), int64(2), Vector2D{X: 1, Y: 2}, false},
	}
	for _, test := range tests {
		v, err := new(test.a, test.b)
		if err != nil && !test.err {
			t.Errorf("vector2d(%t, %t) returned error %v", test.a, test.b, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("vector2d(%t, %t) returned no error", test.a, test.b)
			continue
		}
		if v != test.v {
			t.Errorf("vector2d(%t, %t) returned %v, want %v", test.a, test.b, v, test.v)
			continue
		}
	}
}

func Test_vector(t *testing.T) {
	testNewVector(t, vector2d)
}

func TestNew(t *testing.T) {
	testNewVector(t, New)
}

func TestFromAngle(t *testing.T) {
	tests := []struct {
		params []interface{}
		v      Vector2D
		err    bool
	}{
		{[]interface{}{0}, Vector2D{1, 0}, false},
		{[]interface{}{math.Pi}, Vector2D{-1, 0}, false},
		{[]interface{}{math.Pi / 2}, Vector2D{0, 1}, false},
		{[]interface{}{math.Pi / 4}, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}, false},
		{[]interface{}{math.Pi / 4, 2}, Vector2D{math.Sqrt2, math.Sqrt2}, false},
		{[]interface{}{math.Pi / 4, 2, 3}, Vector2D{}, true},
		{[]interface{}{math.Pi / 4, "2"}, Vector2D{}, true},
		{[]interface{}{"0"}, Vector2D{}, true},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		v, err := FromAngle(test.params[0], test.params[1:]...)
		if err != nil && !test.err {
			t.Errorf("FromAngle(%t) returned error %v", test.params[0], err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("FromAngle(%t) returned no error", test.params[0])
			continue
		}
		if !cmp.Equal(v, test.v, opt) {
			t.Errorf("FromAngle(%t) returned %v, want %v,%T,%T", test.params[0], v, test.v, v.X, test.v.X)
			continue
		}
	}
}

func checkMagErr(t *testing.T, v Vector2D, mag float32) {
	opt := getComparer(.00001)
	Mag := v.Mag()
	if !cmp.Equal(Mag, mag, opt) {
		t.Errorf("Magnitude(%v) returned %v, want %v", v, Mag, mag)
	}
}

func TestRandom(t *testing.T) {

	defaultMag := float32(1)
	// 1
	v1, err := Random()
	if err != nil {
		t.Errorf("Random() returned error %v", err)
	}
	checkMagErr(t, v1, defaultMag)

	v2, err := Random()
	if err != nil {
		t.Errorf("Random() returned error %v", err)
	}
	checkMagErr(t, v2, defaultMag)

	if v1 == v2 {
		t.Errorf("Random() returned %v, want different", v1)
	}

	// 2
	_, err = Random(1, 2, 3)
	if err == nil {
		t.Errorf("Random(1,2,3) returned no error")
	}

	// 3
	m4 := 1
	v4, err := Random(m4)
	if err != nil {
		t.Errorf("Random(1) returned error %v", err)
	}
	checkMagErr(t, v4, float32(m4))

	// 4
	m5 := 3.4
	v5, err := Random(m5)
	if err != nil {
		t.Errorf("Random(1) returned error %v", err)
	}
	checkMagErr(t, v5, float32(m5))

	// 5
	_, err = Random("2")
	if err == nil {
		t.Errorf("Random(2) returned no error")
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		v Vector2D
		s string
	}{
		{Vector2D{1, 1}, "{X: 1, Y: 1}"},
		{Vector2D{-1, -1}, "{X: -1, Y: -1}"},
		{Vector2D{0, 0}, "{X: 0, Y: 0}"},
		{Vector2D{1, 0}, "{X: 1, Y: 0}"},
		{Vector2D{1.8, 2.6}, "{X: 1.8, Y: 2.6}"},
	}
	for _, test := range tests {
		if test.v.String() != test.s {
			t.Errorf("String(%v) returned %v, want %v", test.v, test.v.String(), test.s)
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		v1, v2     Vector2D
		tolerance  []interface{}
		equal, err bool
	}{
		{Vector2D{1, 2}, Vector2D{1, 2}, []interface{}{}, true, false},
		{Vector2D{2, 2}, Vector2D{2, 2}, []interface{}{0}, true, false},
		{Vector2D{3, 2}, Vector2D{3, 2}, []interface{}{"0"}, false, true},
		{Vector2D{3, 2}, Vector2D{2, 2}, []interface{}{}, false, false},
		{Vector2D{3, 2}, Vector2D{2, 2}, []interface{}{1}, true, false},
		{Vector2D{3, 2}, Vector2D{2, 2}, []interface{}{1, 1}, false, true},
	}
	for _, test := range tests {
		equal, err := test.v1.Equal(test.v2, test.tolerance...)
		if err != nil && !test.err {
			t.Errorf("Equal(%v, %v) returned error %v", test.v1, test.v2, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("Equal(%v, %v) returned no error", test.v1, test.v2)
			continue
		}
		if equal != test.equal {
			t.Errorf("Equal(%v, %v) returned %v, want %v", test.v1, test.v2, equal, test.equal)
			continue
		}
	}
}

func testVecCopy(t *testing.T, copy func(Vector2D) Vector2D) {
	tests := []struct {
		v Vector2D
	}{
		{Vector2D{1, 2}},
		{Vector2D{-1, -2}},
		{Vector2D{0, 0}},
	}
	for _, test := range tests {
		v := copy(test.v)
		if v != test.v {
			t.Errorf("Copy(%v) returned %v, want %v", test.v, v, test.v)
		}
		v.X = 2
		if v == test.v {
			t.Errorf("Copy(%v) did not copy", test.v)
		}
	}
}

func TestVecCopy(t *testing.T) {
	copy := func(v Vector2D) Vector2D {
		return v.Copy()
	}
	testVecCopy(t, copy)
}

func TestCopyVec(t *testing.T) {
	copy := func(v Vector2D) Vector2D {
		return Copy(v)
	}
	testVecCopy(t, copy)
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		v    Vector2D
		magn float32
	}{
		{Vector2D{1, 0}, 1},
		{Vector2D{-1, 0}, 1},
		{Vector2D{0, 1}, 1},
		{Vector2D{0, -1}, 1},
		{Vector2D{1, 1}, math.Sqrt2},
		{Vector2D{0, 0}, 0},
		{Vector2D{3, 4}, 5},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		magn := test.v.Mag()
		if !cmp.Equal(magn, test.magn, opt) {
			t.Errorf("Magnitude(%v) returned %v, want %v", test.v, magn, test.magn)
			continue
		}
	}
}

func TestMagnitudeSqr(t *testing.T) {
	tests := []struct {
		v    Vector2D
		magn float32
	}{
		{Vector2D{1, 0}, 1},
		{Vector2D{-1, 0}, 1},
		{Vector2D{0, 1}, 1},
		{Vector2D{0, -1}, 1},
		{Vector2D{1, 1}, 2},
		{Vector2D{0, 0}, 0},
		{Vector2D{3, 4}, 25},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		magn := test.v.MagSq()
		if !cmp.Equal(magn, test.magn, opt) {
			t.Errorf("MagnitudeSqr(%v) returned %v, want %v", test.v, magn, test.magn)
			continue
		}
	}
}

func TestHeading(t *testing.T) {
	tests := []struct {
		v     Vector2D
		angle float32
	}{
		{Vector2D{1, 0}, 0},
		{Vector2D{-1, 0}, math.Pi},
		{Vector2D{0, 1}, math.Pi / 2},
		{Vector2D{0, -1}, -math.Pi / 2},
		{Vector2D{1, 1}, math.Pi / 4},
		{Vector2D{0, 0}, 0},
		{Vector2D{3, 4}, float32(math.Atan2(4, 3))},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		angle := test.v.Heading()
		if !cmp.Equal(angle, test.angle, opt) {
			t.Errorf("Heading(%v) returned %v, want %v", test.v, angle, test.angle)
			continue
		}
	}
}

func TestUnit(t *testing.T) {
	tests := []struct {
		v    Vector2D
		norm Vector2D
	}{
		{Vector2D{1, 0}, Vector2D{1, 0}},
		{Vector2D{-1, 0}, Vector2D{-1, 0}},
		{Vector2D{0, 2}, Vector2D{0, 1}},
		{Vector2D{0, -2}, Vector2D{0, -1}},
		{Vector2D{1, 1}, Vector2D{1 / math.Sqrt2, 1 / math.Sqrt2}},
		{Vector2D{0, 0}, Vector2D{0, 0}},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		v := Unit(test.v)
		if !cmp.Equal(v, test.norm, opt) {
			t.Errorf("Normalize(%v) returned %v, want %v", test.v, test.v, test.norm)
			continue
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		v    Vector2D
		norm Vector2D
	}{
		{Vector2D{1, 0}, Vector2D{1, 0}},
		{Vector2D{-1, 0}, Vector2D{-1, 0}},
		{Vector2D{0, 2}, Vector2D{0, 1}},
		{Vector2D{0, -2}, Vector2D{0, -1}},
		{Vector2D{1, 1}, Vector2D{1 / math.Sqrt2, 1 / math.Sqrt2}},
		{Vector2D{0, 0}, Vector2D{0, 0}},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		test.v.Normalize()
		if !cmp.Equal(test.v, test.norm, opt) {
			t.Errorf("Normalize(%v) returned %v, want %v", test.v, test.v, test.norm)
			continue
		}
	}
}

func TestResize(t *testing.T) {
	opt := getComparer(.00001)
	tests := []struct {
		v   Vector2D
		m   interface{}
		s   Vector2D
		err bool
	}{
		{Vector2D{0, 0}, 2, Vector2D{0, 0}, false},
		{Vector2D{1, 0}, 2, Vector2D{2, 0}, false},
		{Vector2D{-1, 0}, 2, Vector2D{-2, 0}, false},
		{Vector2D{3, 4}, 10, Vector2D{6, 8}, false},
		{Vector2D{1, 4}, "1.2", Vector2D{1, 4}, true},
	}
	for _, test := range tests {
		err := test.v.Resize(test.m)
		if err != nil && !test.err {
			t.Errorf("Resize(%t) returned error %v", test.m, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("Resize(%t) returned no error", test.m)
			continue
		}
		if !cmp.Equal(test.s, test.v, opt) {
			t.Errorf("Resize(%t) returned %v, want %v", test.m, test.v, test.s)
			continue
		}
	}
}

func TestVecAdd(t *testing.T) {
	opt := getComparer(.00001)
	tests := []struct {
		v1 Vector2D
		v2 Vector2D
		v3 Vector2D
	}{
		{Vector2D{1, 1}, Vector2D{1, 1}, Vector2D{2, 2}},
		{Vector2D{1, 1}, Vector2D{-1, -1}, Vector2D{0, 0}},
		{Vector2D{1, 1}, Vector2D{0, 0}, Vector2D{1, 1}},
		{Vector2D{1, 1}, Vector2D{1, 0}, Vector2D{2, 1}},
		{Vector2D{1, 1}, Vector2D{0, 1}, Vector2D{1, 2}},
	}
	for _, test := range tests {
		test.v1.Add(test.v2)
		if !cmp.Equal(test.v1, test.v3, opt) {
			t.Errorf("Add(%v, %v) returned %v, want %v", test.v1, test.v2, test.v1, test.v3)
			continue
		}
	}
}

func TestAddVec(t *testing.T) {
	opt := getComparer(.00001)
	tests := []struct {
		v1 Vector2D
		v2 Vector2D
		v3 Vector2D
	}{
		{Vector2D{1, 1}, Vector2D{1, 1}, Vector2D{2, 2}},
		{Vector2D{1, 1}, Vector2D{-1, -1}, Vector2D{0, 0}},
		{Vector2D{1, 1}, Vector2D{0, 0}, Vector2D{1, 1}},
		{Vector2D{1, 1}, Vector2D{1, 0}, Vector2D{2, 1}},
		{Vector2D{1, 1}, Vector2D{0, 1}, Vector2D{1, 2}},
	}
	for _, test := range tests {
		v := Add(test.v1, test.v2)
		if !cmp.Equal(v, test.v3, opt) {
			t.Errorf("Add(%v, %v) returned %v, want %v", test.v1, test.v2, v, test.v3)
			continue
		}
	}
}

func TestVecSub(t *testing.T) {
	opt := getComparer(.00001)
	tests := []struct {
		v1 Vector2D
		v2 Vector2D
		v3 Vector2D
	}{
		{Vector2D{1, 1}, Vector2D{1, 1}, Vector2D{0, 0}},
		{Vector2D{1, 1}, Vector2D{-1, -1}, Vector2D{2, 2}},
		{Vector2D{1, 1}, Vector2D{0, 0}, Vector2D{1, 1}},
		{Vector2D{1, 1}, Vector2D{1, 0}, Vector2D{0, 1}},
		{Vector2D{1, 1}, Vector2D{0, 1}, Vector2D{1, 0}},
	}
	for _, test := range tests {
		test.v1.Sub(test.v2)
		if !cmp.Equal(test.v1, test.v3, opt) {
			t.Errorf("Sub(%v, %v) returned %v, want %v", test.v1, test.v2, test.v1, test.v3)
			continue
		}
	}
}

func TestSubVec(t *testing.T) {
	opt := getComparer(.00001)
	tests := []struct {
		v1 Vector2D
		v2 Vector2D
		v3 Vector2D
	}{
		{Vector2D{1, 1}, Vector2D{1, 1}, Vector2D{0, 0}},
		{Vector2D{1, 1}, Vector2D{-1, -1}, Vector2D{2, 2}},
		{Vector2D{1, 1}, Vector2D{0, 0}, Vector2D{1, 1}},
		{Vector2D{1, 1}, Vector2D{1, 0}, Vector2D{0, 1}},
		{Vector2D{1, 1}, Vector2D{0, 1}, Vector2D{1, 0}},
	}
	for _, test := range tests {
		v := Sub(test.v1, test.v2)
		if !cmp.Equal(v, test.v3, opt) {
			t.Errorf("Sub(%v, %v) returned %v, want %v", test.v1, test.v2, v, test.v3)
			continue
		}
	}
}

func TestMult(t *testing.T) {
	tests := []struct {
		v   Vector2D
		m   interface{}
		s   Vector2D
		err bool
	}{
		{Vector2D{0, 0}, 2, Vector2D{0, 0}, false},
		{Vector2D{1, 0}, 2, Vector2D{2, 0}, false},
		{Vector2D{-1, 0}, 2, Vector2D{-2, 0}, false},
		{Vector2D{1, 4}, 1.2, Vector2D{1.2, 4.8}, false},
		{Vector2D{1, 4}, 0, Vector2D{0, 0}, false},
		{Vector2D{1, 4}, "1.2", Vector2D{1, 4}, true},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		err := test.v.Mult(test.m)
		if err != nil && !test.err {
			t.Errorf("Scale(%t) returned error %v", test.m, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("Scale(%t) returned no error", test.m)
			continue
		}
		if !cmp.Equal(test.s, test.v, opt) {
			t.Errorf("Scale(%t) returned %v, want %v", test.m, test.v, test.s)
			continue
		}
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		v   Vector2D
		m   interface{}
		s   Vector2D
		err bool
	}{
		{Vector2D{0, 0}, 2, Vector2D{0, 0}, false},
		{Vector2D{2, 3}, 0, Vector2D{2, 3}, true},
		{Vector2D{1, 0}, 2, Vector2D{.5, 0}, false},
		{Vector2D{-1, 0}, 2, Vector2D{-.5, 0}, false},
		{Vector2D{1, 4}, 1.2, Vector2D{1 / 1.2, 4 / 1.2}, false},
		{Vector2D{1, 4}, "1.2", Vector2D{1, 4}, true},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		err := test.v.Div(test.m)
		if err != nil && !test.err {
			t.Errorf("Scale(%t) returned error %v", test.m, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("Scale(%t) returned no error", test.m)
			continue
		}
		if !cmp.Equal(test.s, test.v, opt) {
			t.Errorf("Scale(%t) returned %v, want %v", test.m, test.v, test.s)
			continue
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		v   Vector2D
		r   interface{}
		s   Vector2D
		err bool
	}{
		{Vector2D{1, 0}, 0, Vector2D{1, 0}, false},
		{Vector2D{1, 0}, math.Pi, Vector2D{-1, 0}, false},
		{Vector2D{1, 0}, math.Pi / 2, Vector2D{0, 1}, false},
		{Vector2D{1, 0}, -math.Pi / 2, Vector2D{0, -1}, false},
		{Vector2D{1, 0}, math.Pi / 4, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}, false},
		{Vector2D{1, 0}, -math.Pi / 4, Vector2D{math.Sqrt2 / 2, -math.Sqrt2 / 2}, false},
		{Vector2D{1, 0}, "0", Vector2D{1, 0}, true},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		v := test.v
		err := v.Rotate(test.r)
		if err != nil && !test.err {
			t.Errorf("Rotate(%v, %v) returned error %v, want no error", v, test.r, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("Rotate(%v, %v) returned no error, want error", v, test.r)
			continue
		}
		if !cmp.Equal(v, test.s, opt) {
			t.Errorf("Rotate(%v, %v) returned %v, want %v", v, test.r, v, test.s)
		}
	}
}

func TestSetHeading(t *testing.T) {
	tests := []struct {
		v   Vector2D
		h   interface{}
		s   Vector2D
		err bool
	}{
		{Vector2D{1, 0}, 0, Vector2D{1, 0}, false},
		{Vector2D{1, 0}, math.Pi, Vector2D{-1, 0}, false},
		{Vector2D{1, 0}, math.Pi / 2, Vector2D{0, 1}, false},
		{Vector2D{1, 0}, -math.Pi / 2, Vector2D{0, -1}, false},
		{Vector2D{1, 0}, math.Pi / 4, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}, false},
		{Vector2D{1, 0}, -math.Pi / 4, Vector2D{math.Sqrt2 / 2, -math.Sqrt2 / 2}, false},
		{Vector2D{1, 0}, "0", Vector2D{1, 0}, true},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		v := test.v
		err := v.SetHeading(test.h)
		if err != nil && !test.err {
			t.Errorf("SetHeading(%v, %v) returned error %v, want no error", v, test.h, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("SetHeading(%v, %v) returned no error, want error", v, test.h)
			continue
		}
		if !cmp.Equal(v, test.s, opt) {
			t.Errorf("SetHeading(%v, %v) returned %v, want %v", v, test.h, v, test.s)
		}
	}
}

func TestDist(t *testing.T) {
	tests := []struct {
		v1, v2 Vector2D
		d      float32
	}{
		{Vector2D{0, 0}, Vector2D{0, 0}, 0},
		{Vector2D{1, 0}, Vector2D{0, 0}, 1},
		{Vector2D{0, 1}, Vector2D{0, 0}, 1},
		{Vector2D{1, 1}, Vector2D{0, 0}, math.Sqrt2},
		{Vector2D{3, 4}, Vector2D{4, 3}, math.Sqrt2},
		{Vector2D{1, 0}, Vector2D{0, 1}, math.Sqrt2},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		d := test.v1.Dist(test.v2)
		if !cmp.Equal(d, test.d, opt) {
			t.Errorf("Dist(%v, %v) returned %v, want %v", test.v1, test.v2, d, test.d)
		}
	}
}

func testDot(t *testing.T, dot func(Vector2D, Vector2D) float32) {
	tests := []struct {
		v1, v2 Vector2D
		d      float32
	}{
		{Vector2D{0, 0}, Vector2D{0, 0}, 0},
		{Vector2D{1, 0}, Vector2D{0, 0}, 0},
		{Vector2D{0, 1}, Vector2D{1, 0}, 0},
		{Vector2D{3, 4}, Vector2D{4, 3}, 24},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		d := dot(test.v1, test.v2)
		if !cmp.Equal(d, test.d, opt) {
			t.Errorf("Dot(%v, %v) returned %v, want %v", test.v1, test.v2, d, test.d)
		}
	}
}

func TestVecDot(t *testing.T) {
	dot := func(v1, v2 Vector2D) float32 {
		return v1.Dot(v2)
	}
	testDot(t, dot)
}

func TestDotVec(t *testing.T) {
	dot := func(v1, v2 Vector2D) float32 {
		return Dot(v1, v2)
	}
	testDot(t, dot)
}

func testCross(t *testing.T, cross func(v1, v2 Vector2D) float32) {
	tests := []struct {
		v1, v2 Vector2D
		d      float32
	}{
		{Vector2D{0, 0}, Vector2D{0, 0}, 0},
		{Vector2D{1, 0}, Vector2D{0, 0}, 0},
		{Vector2D{0, 1}, Vector2D{1, 0}, -1},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		d := cross(test.v1, test.v2)
		if !cmp.Equal(d, test.d, opt) {
			t.Errorf("Cross(%v, %v) returned %v, want %v", test.v1, test.v2, d, test.d)
		}
	}
}

func TestVecCross(t *testing.T) {
	cross := func(v1, v2 Vector2D) float32 {
		return v1.Cross(v2)
	}
	testCross(t, cross)
}

func TestCrossVec(t *testing.T) {
	cross := func(v1, v2 Vector2D) float32 {
		return Cross(v1, v2)
	}
	testCross(t, cross)
}

func testAngleBetween(t *testing.T, angleB func(v1, v2 Vector2D) (float32, error)) {
	tests := []struct {
		v1, v2 Vector2D
		a      float32
		err    bool
	}{
		{Vector2D{0, 0}, Vector2D{0, 0}, 0, true},
		{Vector2D{1, 0}, Vector2D{0, 0}, 0, true},
		{Vector2D{0, 0}, Vector2D{1, 1}, 0, true},
		{Vector2D{1, 0}, Vector2D{0, 1}, math.Pi / 2, false},
		{Vector2D{0, 1}, Vector2D{1, 0}, -math.Pi / 2, false},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		a, err := angleB(test.v1, test.v2)
		if err != nil && !test.err {
			t.Errorf("AngleBetween(%v, %v) returned error %v, want no error", test.v1, test.v2, err)
			continue
		}
		if err == nil && test.err {
			t.Errorf("AngleBetween(%v, %v) returned no error, want error", test.v1, test.v2)
			continue
		}
		if !cmp.Equal(a, test.a, opt) {
			t.Errorf("AngleBetween(%v, %v) returned %v, want %v", test.v1, test.v2, a, test.a)
		}
	}
}

func TestVecAngleBetween(t *testing.T) {
	angleB := func(v1, v2 Vector2D) (float32, error) {
		return v1.AngleBetween(v2)
	}
	testAngleBetween(t, angleB)
}
