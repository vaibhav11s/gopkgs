package vector2d

import (
	"math"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testNewVector(t *testing.T, new func(a float32, b float32) Vector2D) {
	var tests []struct {
		a, b float32
		v    Vector2D
	}
	for i := 0; i < 10; i++ {
		a := rand.Float32()*100 - 50
		b := rand.Float32()*100 - 50
		v := new(a, b)
		tests = append(tests, struct {
			a, b float32
			v    Vector2D
		}{a, b, v})
	}
	for _, test := range tests {
		v := new(test.a, test.b)
		if v != test.v {
			t.Errorf("New(%f, %f) = %v, want %v", test.a, test.b, v, test.v)
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
		params []float32
		v      Vector2D
	}{
		{[]float32{0}, Vector2D{1, 0}},
		{[]float32{math.Pi}, Vector2D{-1, 0}},
		{[]float32{math.Pi / 2}, Vector2D{0, 1}},
		{[]float32{math.Pi / 4}, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}},
		{[]float32{math.Pi / 4, 2}, Vector2D{math.Sqrt2, math.Sqrt2}},
		{[]float32{math.Pi / 4, 2, 3}, Vector2D{math.Sqrt2, math.Sqrt2}},
	}
	for _, test := range tests {
		v := FromAngle(test.params[0], test.params[1:]...)
		if !v.Equal(test.v, .00001) {
			t.Errorf("FromAngle(%v) = %v, want %v", test.params, v, test.v)
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
	MAG := float32(1)

	// 1
	v1 := Random()
	checkMagErr(t, v1, MAG)

	// 2
	v2 := Random()
	checkMagErr(t, v2, MAG)

	if v1 == v2 {
		t.Errorf("Random() returned %v, want different", v1)
	}

	// 3
	m3 := float32(1)
	v3 := Random(m3)
	checkMagErr(t, v3, m3)

	// 4
	m4 := float32(3.4)
	v4 := Random(m4)
	checkMagErr(t, v4, m4)

	// 5
	v5 := Random(m4, m3)
	checkMagErr(t, v5, m4)
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
		v1, v2    Vector2D
		tolerance []float32
		equal     bool
	}{
		{Vector2D{1, 2}, Vector2D{1, 2}, []float32{}, true},
		{Vector2D{2, 2}, Vector2D{2, 2}, []float32{0}, true},
		{Vector2D{3, 2}, Vector2D{2, 2}, []float32{}, false},
		{Vector2D{3, 2}, Vector2D{3, 3}, []float32{}, false},
		{Vector2D{3, 2}, Vector2D{2, 2}, []float32{1}, true},
		{Vector2D{3, 2}, Vector2D{2, 2}, []float32{1, 1}, true},
		{Vector2D{3, 2}, Vector2D{2, 2}, []float32{0, 1}, false},
	}
	for _, test := range tests {
		equal := test.v1.Equal(test.v2, test.tolerance...)
		if equal != test.equal {
			t.Errorf("Equal(%v, %v), %v returned %v, want %v", test.v1, test.v2, test.tolerance, equal, test.equal)
			continue
		}
	}
	v1 := Vector2D{1, 2}
	v2 := Vector2D{1, 2}
	V1 := &v1
	V2 := &v2
	if !v1.Equal(v2) {
		t.Errorf("Equal(%v, %v) returned %v, want %v", v1, v2, v1.Equal(v2), true)
	}
	if !v1.Equal(*V2) {
		t.Errorf("Equal(%v, *&%v) returned %v, want %v", v1, V2, v1.Equal(*V2), true)
	}
	if !V1.Equal(v2) {
		t.Errorf("Equal(&%v, %v) returned %v, want %v", V1, v2, V1.Equal(v2), true)
	}
	if !V1.Equal(*V2) {
		t.Errorf("Equal(&%v, *&%v) returned %v, want %v", V1, V2, V1.Equal(*V2), true)
	}
}

func testVecCopy(t *testing.T, copy func(*Vector2D) *Vector2D) {
	tests := []struct {
		v *Vector2D
	}{
		{&Vector2D{1, 2}},
		{&Vector2D{-1, -2}},
		{&Vector2D{3, 4}},
		{&Vector2D{-3, 8}},
	}
	for _, test := range tests {
		v := copy(test.v)
		if !v.Equal(*test.v) {
			t.Errorf("Copy(%v) returned %v, want %v", test.v, v, test.v)
		}
		v.X = 32
		if v.Equal(*test.v) {
			t.Errorf("Changing values of Copy(%v) did change original", test.v)
		}
		test.v.X = 21
		if v.Equal(*test.v) {
			t.Errorf("Changing values of original did change Copy(%v)", test.v)
		}

	}
}

func TestVecCopy(t *testing.T) {
	copy := func(v *Vector2D) *Vector2D {
		v1 := v.Copy()
		return &v1
	}
	testVecCopy(t, copy)
}

func TestCopyVec(t *testing.T) {
	copy := func(v *Vector2D) *Vector2D {
		v1 := Copy(*v)
		return &v1
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
		{Vector2D{-3, 4}, 5},
		{Vector2D{0.6, 0.8}, 1},
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
	for _, test := range tests {
		v := Unit(test.v)
		if !v.Equal(test.norm, .00001) {
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
		{Vector2D{3, 4}, Vector2D{3.0 / 5, 4.0 / 5}},
	}
	for _, test := range tests {
		test.v.Normalize()
		if !test.v.Equal(test.norm, .00001) {
			t.Errorf("Normalize(%v) returned %v, want %v", test.v, test.v, test.norm)
			continue
		}
	}
}

func TestResize(t *testing.T) {
	tests := []struct {
		v Vector2D
		m float32
		s Vector2D
	}{
		{Vector2D{0, 0}, 2, Vector2D{0, 0}},
		{Vector2D{1, 0}, 2, Vector2D{2, 0}},
		{Vector2D{-1, 0}, 2, Vector2D{-2, 0}},
		{Vector2D{3, 4}, 10, Vector2D{6, 8}},
	}
	for _, test := range tests {
		test.v.Resize(test.m)
		if !test.s.Equal(test.v, .00001) {
			t.Errorf("Resize(%v, %v) returned %v, want %v", test.v, test.m, test.v, test.s)
			continue
		}
	}
}

func TestVecAdd(t *testing.T) {
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
		v1 := test.v1.Copy()
		v2 := test.v2.Copy()
		test.v1.Add(test.v2)
		if !test.v1.Equal(test.v3, .00001) {
			t.Errorf("Add(%v, %v) returned %v, want %v", test.v1, test.v2, test.v1, test.v3)
			continue
		}
		if !v2.Equal(test.v2) {
			t.Errorf("Add(%v, %v) changed v2 to %v, want %v", v1, v2, test.v2, v2)
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
		v1 := test.v1.Copy()
		v2 := test.v2.Copy()
		v := Add(test.v1, test.v2)
		if !cmp.Equal(v, test.v3, opt) {
			t.Errorf("Add(%v, %v) returned %v, want %v", test.v1, test.v2, v, test.v3)
			continue
		}
		if !v1.Equal(test.v1) {
			t.Errorf("Add(%v, %v) changed v1 to %v, want %v", v1, v2, test.v1, v1)
			continue
		}
		if !v2.Equal(test.v2) {
			t.Errorf("Add(%v, %v) changed v2 to %v, want %v", v1, v2, test.v2, v2)
			continue
		}
	}
}

func TestVecSub(t *testing.T) {
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
		v1 := test.v1.Copy()
		v2 := test.v2.Copy()
		test.v1.Sub(test.v2)
		if !test.v1.Equal(test.v3, .00001) {
			t.Errorf("Sub(%v, %v) returned %v, want %v", test.v1, test.v2, test.v1, test.v3)
			continue
		}
		if !v2.Equal(test.v2) {
			t.Errorf("Sub(%v, %v) changed v2 to %v, want %v", v1, v2, test.v2, v2)
			continue
		}
	}
}

func TestSubVec(t *testing.T) {
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
		v1 := test.v1.Copy()
		v2 := test.v2.Copy()
		v := Sub(test.v1, test.v2)
		if !v.Equal(test.v3, .00001) {
			t.Errorf("Sub(%v, %v) returned %v, want %v", test.v1, test.v2, v, test.v3)
			continue
		}
		if !v1.Equal(test.v1) {
			t.Errorf("Sub(%v, %v) changed v1 to %v, want %v", v1, v2, test.v1, v1)
			continue
		}
		if !v2.Equal(test.v2) {
			t.Errorf("Sub(%v, %v) changed v2 to %v, want %v", v1, v2, test.v2, v2)
			continue
		}
	}
}

func TestMult(t *testing.T) {
	tests := []struct {
		v Vector2D
		m float32
		s Vector2D
	}{
		{Vector2D{0, 0}, 2, Vector2D{0, 0}},
		{Vector2D{1, 0}, 2, Vector2D{2, 0}},
		{Vector2D{-1, 0}, 2, Vector2D{-2, 0}},
		{Vector2D{1, 4}, 1.2, Vector2D{1.2, 4.8}},
		{Vector2D{1, 4}, 0, Vector2D{0, 0}},
	}
	for _, test := range tests {
		test.v.Mult(test.m)
		if !test.s.Equal(test.v, .00001) {
			t.Errorf("Mult(%v, %v) returned %v, want %v", test.v, test.m, test.v, test.s)
			continue
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		v Vector2D
		r float32
		s Vector2D
	}{
		{Vector2D{1, 0}, 0, Vector2D{1, 0}},
		{Vector2D{1, 0}, math.Pi, Vector2D{-1, 0}},
		{Vector2D{1, 0}, math.Pi / 2, Vector2D{0, 1}},
		{Vector2D{1, 0}, -math.Pi / 2, Vector2D{0, -1}},
		{Vector2D{1, 0}, math.Pi / 4, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}},
		{Vector2D{1, 0}, -math.Pi / 4, Vector2D{math.Sqrt2 / 2, -math.Sqrt2 / 2}},
	}
	for _, test := range tests {
		v := test.v
		v.Rotate(test.r)
		if !v.Equal(test.s, .00001) {
			t.Errorf("Rotate(%v, %v) returned %v, want %v", v, test.r, v, test.s)
		}
	}
}

func TestSetHeading(t *testing.T) {
	tests := []struct {
		v Vector2D
		h float32
		s Vector2D
	}{
		{Vector2D{1, 0}, 0, Vector2D{1, 0}},
		{Vector2D{1, 0}, math.Pi, Vector2D{-1, 0}},
		{Vector2D{1, 0}, math.Pi / 2, Vector2D{0, 1}},
		{Vector2D{1, 0}, -math.Pi / 2, Vector2D{0, -1}},
		{Vector2D{1, 0}, math.Pi / 4, Vector2D{math.Sqrt2 / 2, math.Sqrt2 / 2}},
		{Vector2D{1, 0}, -math.Pi / 4, Vector2D{math.Sqrt2 / 2, -math.Sqrt2 / 2}},
	}
	for _, test := range tests {
		v := test.v
		v.SetHeading(test.h)
		if !v.Equal(test.s, .00001) {
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

func testAngleBetween(t *testing.T, angleB func(v1, v2 Vector2D) float32) {
	NaN := float32(math.NaN())
	tests := []struct {
		v1, v2 Vector2D
		a      float32
	}{
		{Vector2D{0, 0}, Vector2D{0, 0}, NaN},
		{Vector2D{1, 0}, Vector2D{0, 0}, NaN},
		{Vector2D{0, 0}, Vector2D{1, 1}, NaN},
		{Vector2D{1, 0}, Vector2D{0, 1}, math.Pi / 2},
		{Vector2D{0, 1}, Vector2D{1, 0}, -math.Pi / 2},
		{Vector2D{0, 1}, Vector2D{0, 1}, 0},
		{Vector2D{0, 2}, Vector2D{0, 12}, 0},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		a := angleB(test.v1, test.v2)
		if math.IsNaN(float64(a)) && math.IsNaN(float64(test.a)) {
			continue
		}
		if !cmp.Equal(a, test.a, opt) {
			t.Errorf("AngleBetween(%v, %v) returned %v, want %v", test.v1, test.v2, a, test.a)
		}
	}
}

func TestVecAngleBetween(t *testing.T) {
	angleB := func(v1, v2 Vector2D) float32 {
		return v1.AngleBetween(v2)
	}
	testAngleBetween(t, angleB)
}

func TestAngleBetweenVec(t *testing.T) {
	angleB := func(v1, v2 Vector2D) float32 {
		return AngleBetween(v1, v2)
	}
	testAngleBetween(t, angleB)
}
