// Package vector provides a simple 2D vector type on cartesian plane
package vector2d

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Vector2D struct {
	X, Y float32
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Creates a new 2D vector.
// Two dimensional Euclidean vector.
func New(x, y float32) *Vector2D {
	return &Vector2D{x, y}
}

// Make a new 2D vector from an angle
func FromAngle(angle float32, length ...float32) *Vector2D {
	var l float32 = 1
	if len(length) >= 1 {
		l = length[0]
	}
	return &Vector2D{float32(math.Cos(float64(angle)) * float64(l)), float32(math.Sin(float64(angle)) * float64(l))}
}

// Make a new 2D vector from a random angle of length 1 (default) or a given length
func Random(length ...float32) *Vector2D {
	var l float32 = 1
	if len(length) >= 1 {
		l = length[0]
	}
	ang := rand.Float32() * 2 * math.Pi
	return FromAngle(ang, l)
}

// Returns a string representation of the vector
func (v *Vector2D) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", v.X, v.Y)
}

// Checks whether two vectors are equal.
// optional tolerence value can be passed as a parameter to check for equality
// within a tolerance, abs(v.x - v2.x) < tolerance and abs(v.y - v2.y) < tolerance
func (v *Vector2D) Equal(v2 *Vector2D, tolerance ...float32) bool {
	var t float32 = 0
	if len(tolerance) >= 1 {
		t = tolerance[0]
	}
	if math.Abs(float64(v.X-v2.X)) > float64(t) {
		return false
	}
	if math.Abs(float64(v.Y-v2.Y)) > float64(t) {
		return false
	}
	return true
}

// Gets a copy of the vector
func (v *Vector2D) Copy() *Vector2D {
	return &Vector2D{v.X, v.Y}
}

// Calculates the magnitude (length) of the vector and returns the result as a float
// this is simply the equation sqrt(x*x + y*y)
func (v *Vector2D) Mag() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Calculates the squared magnitude of the vector and returns the result as a float
// this is simply the equation (x*x + y*y + z*z)
func (v *Vector2D) MagSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Calculate the angle of rotation for the vector
func (v *Vector2D) Heading() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

// Normalize the vector to length 1 (make it a unit vector).
// Modify + Returns self
func (v *Vector2D) Normalize() *Vector2D {
	m := v.Mag()
	if m == 0 {
		return v
	}
	v.X /= m
	v.Y /= m
	return v
}

// Set the magnitude of this vector to the value used for the len parameter.
// Modify + Returns self
func (v *Vector2D) Resize(mag float32) *Vector2D {
	m := v.Mag()
	if m == 0 {
		return v
	}
	v.X = v.X * mag / m
	v.Y = v.Y * mag / m
	return v
}

// add a vector to the current vector.
// Modify + Returns self
func (v *Vector2D) Add(v2 *Vector2D) *Vector2D {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

// subtract a vector from the current vector.
// Modify + Returns self
func (v *Vector2D) Sub(v2 *Vector2D) *Vector2D {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

// Multiplies the vector by a scalar.
// Modify + Returns self
func (v *Vector2D) Mult(scalar float32) *Vector2D {
	v.X *= scalar
	v.Y *= scalar
	return v
}

// rotate the vector in the direction of the angle.
// Modify + Returns self
func (v *Vector2D) Rotate(angle float32) *Vector2D {
	newHeading := v.Heading() + angle
	m := v.Mag()
	v.X = float32(math.Cos(float64(newHeading))) * m
	v.Y = float32(math.Sin(float64(newHeading))) * m
	return v
}

// Rotate the vector to a specific angle, magnitude remains the same.
// Modify + Returns self
func (v *Vector2D) SetHeading(angle float32) *Vector2D {
	m := v.Mag()
	v.X = float32(math.Cos(float64(angle))) * m
	v.Y = float32(math.Sin(float64(angle))) * m
	return v
}

// Calculates the Euclidean distance between two points
// (considering a point as a vector object)
func (v *Vector2D) Dist(v2 *Vector2D) float32 {
	sV := Sub(v, v2)
	return sV.Mag()
}

// Calculates the dot product with another vector
func (v *Vector2D) Dot(v2 *Vector2D) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Calculates the cross product with another vector
// ~ give the value of the z axis component
// (in 2D space, the cross product is a vector perpendicular to the two input vectors)
func (v *Vector2D) Cross(v2 *Vector2D) float32 {
	return v.X*v2.Y - v.Y*v2.X
}

// Calculates and returns the angle with another vector
// Returns NaN if any vector is a zero vector
func (v *Vector2D) AngleBetween(v2 *Vector2D) float32 {
	m1 := v.Mag()
	m2 := v2.Mag()
	if m1 == 0 || m2 == 0 {
		return float32(math.NaN())
	}
	dotMag := Dot(v, v2) / (m1 * m2)
	angle := math.Acos(math.Min(1, math.Max(-1, float64(dotMag))))
	sign := Cross(v, v2) < 0
	if sign {
		angle = -angle
	}
	return float32(angle)
}

// Gets a copy of the vector
func Copy(v *Vector2D) *Vector2D {
	return &Vector2D{v.X, v.Y}
}

// Gives a unit vector in dirction of the vector
func Unit(v *Vector2D) *Vector2D {
	m := v.Mag()
	if m == 0 {
		return &Vector2D{0, 0}
	}
	return &Vector2D{v.X / m, v.Y / m}
}

// returns the sum of two vectors
func Add(v1, v2 *Vector2D) *Vector2D {
	return &Vector2D{v1.X + v2.X, v1.Y + v2.Y}
}

// returns the difference of two vectors
func Sub(v1, v2 *Vector2D) *Vector2D {
	return &Vector2D{v1.X - v2.X, v1.Y - v2.Y}
}

// Calculates the dot product of two vectors
func Dot(v1, v2 *Vector2D) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Calculates the cross product of two vectors
// ~ give the value of the z axis component
// (in 2D space, the cross product is a vector perpendicular to the two input vectors)
func Cross(v1, v2 *Vector2D) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Calculates and returns the angle between two vectors.
// Returns NaN if any vector is a zero vector
func AngleBetween(v1, v2 *Vector2D) float32 {
	m1 := v1.Mag()
	m2 := v2.Mag()
	if m1 == 0 || m2 == 0 {
		return float32(math.NaN())
	}
	dotMag := Dot(v1, v2) / (m1 * m2)
	angle := math.Acos(math.Min(1, math.Max(-1, float64(dotMag))))
	sign := Cross(v1, v2) < 0
	if sign {
		angle = -angle
	}
	return float32(angle)
}
