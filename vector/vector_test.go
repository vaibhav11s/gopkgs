package vector

import (
	"math"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getComparer(tolerance float64) cmp.Option {
	return cmp.Comparer(func(x, y float32) bool {
		diff := math.Abs(float64(x - y))
		return diff <= tolerance
	})
}

func TestFromAngles(t *testing.T) {
	P4 := float32(math.Pi / 4)
	P2 := float32(math.Pi / 2)
	tests := []struct {
		theta, phi float32
		length     []float32
		want       *Vector
	}{
		{P4, 0.9553166, []float32{}, New(1, 1, 1).Normalize()},
		{P4, 0.9553166, []float32{2}, New(1, 1, 1).Normalize().Mult(2)},
		{rand.Float32() * 4 * P2, 0, []float32{}, New(0, 0, 1)},
		{rand.Float32() * 4 * P2, math.Pi, []float32{}, New(0, 0, -1)},
		//
		{0 * P2, P4, []float32{}, New(1, 0, 1).Normalize()},
		{1 * P2, P4, []float32{}, New(0, 1, 1).Normalize()},
		{2 * P2, P4, []float32{}, New(-1, 0, 1).Normalize()},
		{3 * P2, P4, []float32{}, New(0, -1, 1).Normalize()},
		//
		{0 * P2, 3 * P4, []float32{}, New(1, 0, -1).Normalize()},
		{1 * P2, 3 * P4, []float32{}, New(0, 1, -1).Normalize()},
		{2 * P2, 3 * P4, []float32{}, New(-1, 0, -1).Normalize()},
		{3 * P2, 3 * P4, []float32{}, New(0, -1, -1).Normalize()},
		//
		{1 * P4, P2, []float32{}, New(1, 1, 0).Normalize()},
		{3 * P4, P2, []float32{}, New(-1, 1, 0).Normalize()},
		{5 * P4, P2, []float32{}, New(-1, -1, 0).Normalize()},
		{7 * P4, P2, []float32{}, New(1, -1, 0).Normalize()},
		//
		{1 * P4, 3 * P2, []float32{}, New(-1, -1, 0).Normalize()},
		{3 * P4, 3 * P2, []float32{}, New(1, -1, 0).Normalize()},
		{5 * P4, 3 * P2, []float32{}, New(1, 1, 0).Normalize()},
		{7 * P4, 3 * P2, []float32{}, New(-1, 1, 0).Normalize()},
		//
		{0 * P2, P2, []float32{}, New(1, 0, 0)},
		{1 * P2, P2, []float32{}, New(0, 1, 0)},
		{2 * P2, P2, []float32{}, New(-1, 0, 0)},
		{3 * P2, P2, []float32{}, New(0, -1, 0)},
		//
		{1 * P4, P4, []float32{}, New(1, 1, 1.4142131).Normalize()},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		for i := 0; i < 4; i += 1 {

			if v := FromAngles(test.theta+4*P2*float32(i), test.phi+4*P2*float32(i), test.length...); !cmp.Equal(v, test.want, opt) {
				t.Errorf("FromAngles(%v, %v, %v) = %v, want %v", test.theta, test.phi, test.length, v, test.want)
			}
		}
	}
}

func TestRandom(t *testing.T) {
	opt := getComparer(.00001)
	// Mag = 1
	for i := 0; i < 100; i += 1 {
		v := Random()
		x := v.Copy()
		X := Unit(v)
		if !cmp.Equal(x, X, opt) {
			t.Errorf("got %v, want %v", x, X)
		}
	}
	// Mag = i
	for i := float32(1); i < 100; i += 1 {
		v := Random(i)
		x := v.Copy()
		X := v.Copy().Resize(i)
		if !cmp.Equal(x, X, opt) {
			t.Errorf("got %v, want %v", x, X)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		v    *Vector
		want string
	}{
		{zero(), "{X: 0, Y: 0, Z: 0}"},
		{New(1.232, 0, 64.0), "{X: 1.232, Y: 0, Z: 64}"},
	}
	for _, test := range tests {
		if got := test.v.String(); got != test.want {
			t.Errorf("%v.String() = %v, want %v", test.v, got, test.want)
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b      *Vector
		tolerance []float32
		want      bool
	}{
		{New(1, 2, 3), New(1, 2, 3), []float32{}, true},
		{New(1, 2, 3), New(1, 2.1, 3), []float32{}, false},
		{New(1.1, 2.3, 3.2), New(1.1, 2.3, 3.2), []float32{0.1}, true},
		{New(1.1, 2.4, 3.2), New(1.0, 2.5, 3.2), []float32{0.1}, true},
		{New(1.1, 2.4, 3.2), New(1.0, 2.5, 3.5), []float32{0.1}, false},
		{New(1.4, 2.2, 3.2), New(1.0, 2.21, 3.4), []float32{0.01}, false},
		{New(1.41, 2.2, 3.2), New(1.4, 2.22, 3.19), []float32{0.01}, false},
		{New(1.41, 2.2, 3.2), New(1.4, 2.20, 3.19), []float32{0.01}, true},
	}
	for _, test := range tests {
		if got := test.a.Equal(test.b, test.tolerance...); got != test.want {
			t.Errorf("%v.Equal(%v, %v...) = %v, want %v", test.a, test.b, test.tolerance, got, test.want)
		}
	}
}

func TestVCopy(t *testing.T) {
	tests := []struct {
		v1   *Vector
		want *Vector
	}{
		{New(1, 2, 3), New(1, 2, 3)},
		{New(2, 5, 8), New(2, 5, 8)},
		{New(12, 19, 6), New(12, 19, 6)},
	}
	for _, test := range tests {
		copy := test.v1.Copy()
		if !cmp.Equal(copy, test.want) {
			t.Errorf("%v.Copy() = %v, want %v", test.v1, copy, test.want)
		}
		if copy == test.v1 {
			t.Errorf("%v.Copy() = %v, want a copy", test.v1, copy)
		}
		copy.X = 32
		if cmp.Equal(copy, test.v1) {
			t.Errorf("changing copy also changed original %v", test.v1)
		}
		test.v1.Y = 43
		if cmp.Equal(test.v1, copy) {
			t.Errorf("changing original also changed copy %v", copy)
		}
	}
}

func TestCopyV(t *testing.T) {
	tests := []struct {
		v1   *Vector
		want *Vector
	}{
		{New(1, 2, 3), New(1, 2, 3)},
		{New(2, 5, 8), New(2, 5, 8)},
		{New(12, 19, 6), New(12, 19, 6)},
	}
	for _, test := range tests {
		copy := Copy(test.v1)
		if !cmp.Equal(copy, test.want) {
			t.Errorf("Copy(%v) = %v, want %v", test.v1, copy, test.want)
		}
		if copy == test.v1 {
			t.Errorf("Copy(%v) = %v, want a copy", test.v1, copy)
		}
		copy.X = 32
		if cmp.Equal(copy, test.v1) {
			t.Errorf("changing copy also changed original %v", test.v1)
		}
		test.v1.Y = 43
		if cmp.Equal(test.v1, copy) {
			t.Errorf("changing original also changed copy %v", copy)
		}
	}
}

func TestAssign(t *testing.T) {
	tests := []struct {
		want *Vector
	}{
		{New(1, 2, 3)},
		{New(2, 5, 8)},
		{New(12, 19, 6)},
	}
	for _, test := range tests {
		v := zero()
		v.Assign(test.want)
		if !cmp.Equal(v, test.want) {
			t.Errorf("%v.Assign(%v) = %v, want %v", v, test.want, v, test.want)
		}
		if v == test.want {
			t.Errorf("%v.Assign(%v) = %v, want a copy", v, test.want, v)
		}
		v.X = 32
		if cmp.Equal(v, test.want) {
			t.Errorf("changing copy also changed original %v", test.want)
		}
		test.want.Y = 43
		if cmp.Equal(test.want, v) {
			t.Errorf("changing original also changed copy %v", v)
		}
	}
}

func TestMag(t *testing.T) {
	tests := []struct {
		v    *Vector
		want float32
	}{
		{New(0.26726, 0.53452, 0.80178), 1},
		{New(1, 2, 3), float32(math.Sqrt(14))},
		{New(3, 4, 12), 13},
		{New(2, 5, 8), float32(math.Sqrt(93))},
		{New(12, 19, 6), float32(math.Sqrt(541))},
		{New(312.1511574, 2259.344174, 321.9829745), 2303.4208207624088},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		if got := test.v.Mag(); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.Mag() = %v, want %v", test.v, got, test.want)
		}
	}
}

func TestMagSq(t *testing.T) {
	tests := []struct {
		v    *Vector
		want float32
	}{
		{New(0.26726, 0.53452, 0.80178), 1},
		{New(1, 2, 3), 14},
		{New(3, 4, 12), 169},
		{New(2, 5, 8), 93},
		{New(12, 19, 6), 541},
		{New(12.15, 59.34, 21.92), 4149.3445},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		if got := test.v.MagSq(); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.MagSq() = %v, want %v", test.v, got, test.want)
		}
	}
}

func TestNormalizeUnit(t *testing.T) {
	tests := []struct {
		v    *Vector
		want *Vector
	}{
		{zero(), zero()},
		{New(1, 2, 3), New(0.26726, 0.53452, 0.80178)},
		{New(2, 5, 8), New(float32(2/math.Sqrt(93)), float32(5/math.Sqrt(93)), float32(8/math.Sqrt(93)))},
		{New(3, 4, 12), New(3/13.0, 4/13.0, 12/13.0)},
	}
	opt := getComparer(.00001)
	// Unit
	for _, test := range tests {
		if got := Unit(test.v); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Unit(%v) = %v, want %v", test.v, got, test.want)
		}
	}
	// Normalize
	for _, test := range tests {
		if test.v.Normalize(); !cmp.Equal(test.v, test.want, opt) {
			t.Errorf("%v.Normalize() = %v, want %v", test.v, test.v, test.want)
		}
	}
}

func TestResize(t *testing.T) {
	tests := []struct {
		v    *Vector
		fl   float32
		want *Vector
	}{
		{New(3, 4, 12), 26, New(6, 8, 24)},
		{New(2, 5, 8), 3, New(0.622171, 1.555427, 2.48868)},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		if test.v.Resize(test.fl); !cmp.Equal(test.v, test.want, opt) {
			t.Errorf("%v.Resize(%v) = %v, want %v", test.v, test.fl, test.v, test.want)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		want   *Vector
	}{
		{New(1, 2, 3), New(2, 3, 4), New(3, 5, 7)},
		{New(2, 5, 8), New(12, 19, 6), New(14, 24, 14)},
		{New(8.5, 2, 9.5), New(17.8, 96.0, 3.90), New(26.3, 98.0, 13.4)},
		{New(-68.8, 7.47, 15.9), New(54.5, 48.2, -6.50), New(-14.3, 55.67, 9.4)},
	}
	opt := getComparer(.00001)
	// Add(v1,v2)
	for _, test := range tests {
		if got := Add(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Add(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	// v1.Add(v2)
	for _, test := range tests {
		if test.v1.Add(test.v2); !cmp.Equal(test.v1, test.want, opt) {
			t.Errorf("%v.Add(%v) = %v, want %v", test.v1, test.v2, test.v1, test.want)
		}
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		want   *Vector
	}{
		{New(1, 2, 3), New(2, 3, 4), New(-1, -1, -1)},
		{New(2, 5, 8), New(12, 19, 6), New(-10, -14, 2)},
		{New(3, 4, 12), New(12, 19, 6), New(-9, -15, 6)},
		{New(3, 4, 12), New(-9, -15, 6), New(12, 19, 6)},
		{New(12, 19, 6), New(3, 4, 12), New(9, 15, -6)},
	}
	opt := getComparer(.00001)
	// Sub(v1,v2)
	for _, test := range tests {
		if got := Sub(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Sub(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	// v1.Sub(v2)
	for _, test := range tests {
		if test.v1.Sub(test.v2); !cmp.Equal(test.v1, test.want, opt) {
			t.Errorf("%v.Sub(%v) = %v, want %v", test.v1, test.v2, test.v1, test.want)
		}
	}
}

func TestMult(t *testing.T) {
	tests := []struct {
		v1   *Vector
		fl   float32
		want *Vector
	}{
		{New(1, 2, 3), 2, New(2, 4, 6)},
		{New(2, 5, 8), 12, New(24, 60, 96)},
		{New(3, 4, 12), -2, New(-6, -8, -24)},
		{New(12, 19, 6), 0.5, New(6, 9.5, 3)},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		if test.v1.Mult(test.fl); !cmp.Equal(test.v1, test.want, opt) {
			t.Errorf("%v.Mul(%v) = %v, want %v", test.v1, test.fl, test.v1, test.want)
		}
	}
}

func TestDist(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		want   float32
	}{
		{zero(), zero(), 0},
		{New(1, 2, 3), New(2, 3, 4), 1.732050},
		{New(2, 5, 8), New(12, 19, 6), 17.320508},
	}
	opt := getComparer(.00001)
	// Dist(v1,v2)
	for _, test := range tests {
		if got := Dist(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Dist(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	// v1.Dist(v2)
	for _, test := range tests {
		if got := test.v1.Dist(test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.Dist(%v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		want   float32
	}{
		{New(1, 1, 1), New(1, 1, 1), 3},
		{New(2, 5, 8), New(12, 19, 6), 167},
		{New(1, 1, 1), New(2, -2, 0), 0},
	}
	opt := getComparer(.00001)
	// Dot(v1,v2)
	for _, test := range tests {
		if got := Dot(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Dot(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	// v1.Dot(v2)
	for _, test := range tests {
		if got := test.v1.Dot(test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.Dot(%v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
}

func TestCross(t *testing.T) {
	tests := []struct {
		v1   *Vector
		v2   *Vector
		want *Vector
	}{
		{New(1, 1, 1), New(1, 1, 1), zero()},
		{New(21, 31.2, 12.1), New(4.0, 2.1, 15), New(442.59, -266.6, -80.7)},
		{New(-18, 12.4, -6), New(2.2, -12, 1.2), New(-57.12, 8.4, 188.72)},
		{New(12, 19, 0), New(6, 9, 0), New(0, 0, -6)},
	}
	opt := getComparer(.00001)
	// Cross(v1,v2)
	for _, test := range tests {
		if got := Cross(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Cross(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	// v1.Cross(v2)
	for _, test := range tests {
		if got := test.v1.Cross(test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.Cross(%v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
}

func TestAngle(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		want   float32
	}{
		{New(1, 0, 0), New(0, 1, 0), math.Pi / 2},
		{New(8, 16, 7), New(19, 3, 8), 0.8766778},
		{New(3.2, 2.2, 2.8), New(3.3, 5, 4.2), 0.3135},
	}
	opt := getComparer(.00001)
	// Angle(v1,v2)
	for _, test := range tests {
		if got := Angle(test.v1, test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Angle(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	v := zero()
	v2 := Random()
	if got := Angle(v, v2); !math.IsNaN(float64(got)) {
		t.Errorf("Angle(%v, %v) = %v, want %v", v, v2, got, math.NaN())
	}
	// v1.Angle(v2)
	for _, test := range tests {
		if got := test.v1.Angle(test.v2); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.Angle(%v) = %v, want %v", test.v1, test.v2, got, test.want)
		}
	}
	v = zero()
	v2 = Random()
	if got := v.Angle(v2); !math.IsNaN(float64(got)) {
		t.Errorf("%v.Angle(%v) = %v, want %v", v, v2, got, math.NaN())
	}
}

func TestHeading(t *testing.T) {
	tests := []struct {
		v1         *Vector
		theta, phi float32
	}{
		{New(1, 0, 0), 0, math.Pi / 2},
		{New(0, 1, 0), math.Pi / 2, math.Pi / 2},
		{New(0, 0, 1), 0, 0},
		{New(1, 1, 1), math.Pi / 4, 0.9553166},
		{New(1, -1, 1), -math.Pi / 4, 0.9553166},
		{New(1, 1, -1), math.Pi / 4, math.Pi - 0.9553166},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		th, p := test.v1.Heading()
		if !cmp.Equal(th, test.theta, opt) {
			t.Errorf("%v.Heading()(thetha) = %v, want %v", test.v1, th, test.theta)
		}
		if !cmp.Equal(p, test.phi, opt) {
			t.Errorf("%v.Heading()(phi) = %v, want %v", test.v1, p, test.phi)
		}
	}
	v := zero()
	thetha := float32(0)
	phi := float32(math.NaN())
	th, p := v.Heading()
	if math.IsNaN(float64(phi)) && math.IsNaN(float64(p)) {

	} else {
		t.Errorf("%v.Heading()(phi) = %v, want %v", v, p, phi)
	}
	if !cmp.Equal(th, thetha) {
		t.Errorf("%v.Heading()(thetha) = %v, want %v", v, th, thetha)
	}
}

func TestSetHeading(t *testing.T) {
	test := []struct {
		v1         *Vector
		theta, phi float32
	}{
		{New(1, 0, 0), 0, math.Pi / 2},
		{New(0, 1, 0), math.Pi / 2, math.Pi / 2},
		{New(0, 0, 1), 0, 0},
		{New(1, 1, 1), math.Pi / 4, 0.9553166},
		{New(1, -1, 1), -math.Pi / 4, 0.9553166},
		{New(1, 1, -1), math.Pi / 4, math.Pi - 0.9553166},
	}
	opt := getComparer(.00001)
	for _, test := range test {
		for i := 0; i < 5; i += 1 {
			m := rand.Float32() * 100
			v := Random(m)
			v.SetHeading(test.theta, test.phi)
			test.v1.Resize(m)
			if !cmp.Equal(v, test.v1, opt) {
				t.Errorf("%v.SetHeading(%v, %v) = %v, want %v", v, test.theta, test.phi, v, test.v1)
			}
		}
	}
}

func TestComponent(t *testing.T) {
	tests := []struct {
		v, axis *Vector
		pll, pr *Vector
	}{
		{New(1, 0, 0), New(1, 0, 0), New(1, 0, 0), zero()},
		{New(1, 0, 0), New(0, 1, 0), zero(), New(1, 0, 0)},
		{New(2, 1, 2), zero(), zero(), zero()},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		pl, pr := test.v.Component(test.axis)
		if !cmp.Equal(pr, test.pr, opt) {
			t.Errorf("%v.Component(%v)(perpendicular) = %v, want %v", test.v, test.axis, pr, test.pr)
		}
		if !cmp.Equal(pl, test.pll, opt) {
			t.Errorf("%v.Component(%v)(paraller) = %v, want %v", test.v, test.axis, pl, test.pll)
		}
	}
}

func TestRotatateAlongAxis(t *testing.T) {
	tests := []struct {
		v, axis *Vector
		theta   float32
		want    *Vector
	}{
		{x(), x(), 0, x()},
		{x(), x(), math.Pi / 2, x()},
		{x(), x(), math.Pi, x()},
		{x(), z(), math.Pi / 2, y()},
		{x(), z().Mult(-1), math.Pi / 2, y().Mult(-1)},
		{x(), z(), math.Pi, x().Mult(-1)},
		{x(), z().Mult(-1), math.Pi, x().Mult(-1)},
		{x(), z(), math.Pi / 4, New(1, 1, 0).Normalize()},
		{x(), y(), math.Pi / 2, z().Mult(-1)},
		{x(), y(), math.Pi, x().Mult(-1)},
		{x(), y(), math.Pi / 4, New(1, 0, -1).Normalize()},
		{New(1, 1, 1), New(1, 1, 1), math.Pi / 4, New(1, 1, 1)},
		{New(1, 1, 1), x(), math.Pi / 4, New(1, 0, math.Sqrt2)},
		{New(1, 1, 1), y(), math.Pi / 4, New(math.Sqrt2, 1, 0)},
		{New(1, 1, 1), z(), math.Pi / 4, New(0, math.Sqrt2, 1)},
		{New(1, 1, 1), New(-1, 1, 0), -0.9553166, New(0, 0, float32(math.Sqrt(3)))},
		{New(1, 1, 1), zero(), 0.9553166, New(1, 1, 1)},
		{zero(), Random(), 0.9553166, zero()},
	}
	opt := getComparer(.00001)
	// RotateAlongAxis(v1,ax,th)
	for _, test := range tests {
		v := RotateAlongAxis(test.v, test.axis, test.theta)
		if !cmp.Equal(v, test.want, opt) {
			t.Errorf("RotateAlongAxis(%v, %v, %v) = %v, want %v", test.v, test.axis, test.theta, v, test.want)
		}
	}
	// v1.RotateAlongAxis(ax,th)
	for _, test := range tests {
		test.v.RotateAlongAxis(test.axis, test.theta)
		if !cmp.Equal(test.v, test.want, opt) {
			t.Errorf("%v.RotateAlongAxis(%v, %v) = %v, want %v", test.v, test.axis, test.theta, test.v, test.want)
		}
	}
}

func TestReflectThroughPlane(t *testing.T) {
	tests := []struct {
		v, normal *Vector
		want      *Vector
	}{
		{New(1, 1, 1), New(0, 0, 1), New(1, 1, -1)},
		{zero(), Random(), zero()},
		{New(1, 1, 1), zero(), New(1, 1, 1)},
	}
	opt := getComparer(.00001)
	// ReflectThroughPlane(v,n)
	for _, test := range tests {
		if got := ReflectThroughPlane(test.v, test.normal); !cmp.Equal(got, test.want, opt) {
			t.Errorf("ReflectThroughPlane(%v, %v) = %v, want %v", test.v, test.normal, got, test.want)
		}
	}
	// v.ReflectThroughPlane(n)
	for _, test := range tests {
		if got := test.v.ReflectThroughPlane(test.normal); !cmp.Equal(got, test.want, opt) {
			t.Errorf("%v.ReflectThroughPlane(%v) = %v, want %v", test.v, test.normal, got, test.want)
		}
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		v1, v2 *Vector
		n, i   int
		want   *Vector
	}{
		{zero(), New(1, 1, 1), 2, 0, zero()},
		{zero(), New(1, 1, 1), 2, 1, New(0.5, 0.5, 0.5)},
		{zero(), New(1, 1, 1), 2, 2, New(1, 1, 1)},
	}
	opt := getComparer(.00001)
	for _, test := range tests {
		if got := Lerp(test.v1, test.v2, float32(test.i)/float32(test.n)); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Lerp(%v, %v, %v) = %v, want %v", test.v1, test.v2, float32(test.i)/float32(test.n), got, test.want)
		}
	}
	for _, test := range tests {
		if got := Lerp2(test.v1, test.v2, test.n, test.i); !cmp.Equal(got, test.want, opt) {
			t.Errorf("Lerp(%v, %v, %v) = %v, want %v", test.v1, test.v2, test.i/test.n, got, test.want)
		}
	}
}
