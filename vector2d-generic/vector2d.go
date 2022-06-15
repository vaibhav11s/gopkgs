// Package vector provides a simple 2D vector type on cartesian plane
package vector2d

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

type Vector2D struct {
	X, Y float32
}

var floatType = reflect.TypeOf(float32(0))

func getFloat(unk interface{}) (float32, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("%v is not a float", unk)
	}
	fv := v.Convert(floatType)
	return float32(fv.Float()), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func vector2d(X, Y interface{}) (Vector2D, error) {
	x, err := getFloat(X)
	if err != nil {
		return Vector2D{}, err
	}
	y, err := getFloat(Y)
	if err != nil {
		return Vector2D{}, err
	}
	return Vector2D{x, y}, nil
}

// Creates a new 2D vector.
// Two dimensional Euclidean vector.
func New(x, y interface{}) (Vector2D, error) {
	return vector2d(x, y)
}

// Make a new 2D vector from an angle
func FromAngle(angle interface{}, length ...interface{}) (Vector2D, error) {
	ang, err := getFloat(angle)
	if err != nil {
		return Vector2D{}, err
	}
	if len(length) > 1 {
		return Vector2D{}, fmt.Errorf("too many arguments")
	}
	var l float32 = 1
	if len(length) == 1 {
		l, err = getFloat(length[0])
		if err != nil {
			return Vector2D{}, err
		}
	}
	return Vector2D{float32(math.Cos(float64(ang)) * float64(l)), float32(math.Sin(float64(ang)) * float64(l))}, nil
}

// Make a new 2D vector from a random angle of length 1 (default) or a given length
func Random(length ...interface{}) (Vector2D, error) {
	if len(length) > 1 {
		return Vector2D{}, fmt.Errorf("too many arguments")
	}
	var l float32 = 1
	var err error
	if len(length) == 1 {
		l, err = getFloat(length[0])
		if err != nil {
			return Vector2D{}, err
		}
	}
	ang := rand.Float32() * 2 * math.Pi
	return FromAngle(ang, l)
}

// Returns a string representation of the vector
func (v Vector2D) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", v.X, v.Y)
}

// Checks whether two vectors are equal.
// optional tolerence value can be passed as a parameter to check for equality
// within a tolerance, abs(v.x - v2.x) < tolerance and abs(v.y - v2.y) < tolerance
func (v *Vector2D) Equal(v2 Vector2D, tolerance ...interface{}) (bool, error) {
	if len(tolerance) > 1 {
		return false, fmt.Errorf("too many arguments")
	}
	var t float32 = 0
	var err error
	if len(tolerance) == 1 {
		t, err = getFloat(tolerance[0])
		if err != nil {
			return false, err
		}
	}
	if math.Abs(float64(v.X-v2.X)) > float64(t) {
		return false, nil
	}
	if math.Abs(float64(v.Y-v2.Y)) > float64(t) {
		return false, nil
	}
	return true, nil
}

// Gets a copy of the vector
func (v Vector2D) Copy() Vector2D {
	return Vector2D{v.X, v.Y}
}

// Calculates the magnitude (length) of the vector and returns the result as a float
// this is simply the equation sqrt(x*x + y*y)
func (v Vector2D) Mag() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Calculates the squared magnitude of the vector and returns the result as a float
// this is simply the equation (x*x + y*y + z*z)
func (v Vector2D) MagSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Calculate the angle of rotation for the vector
func (v Vector2D) Heading() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

// Normalize the vector to length 1 (make it a unit vector)
func (v *Vector2D) Normalize() {
	m := v.Mag()
	if m == 0 {
		return
	}
	v.X /= m
	v.Y /= m
}

//  Set the length of this vector to the value used for the len parameter
func (v *Vector2D) Resize(len interface{}) error {
	A, err := getFloat(len)
	if err != nil {
		return err
	}
	m := v.Mag()
	if m == 0 {
		return nil
	}
	v.X = v.X * A / m
	v.Y = v.Y * A / m
	return nil
}

// add a vector to the current vector
func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}

// subtract a vector from the current vector
func (v *Vector2D) Sub(v2 Vector2D) {
	v.X -= v2.X
	v.Y -= v2.Y
}

// Multiplies the vector by a scalar
func (v *Vector2D) Mult(scalar interface{}) error {
	A, err := getFloat(scalar)
	if err != nil {
		return err
	}
	v.X *= A
	v.Y *= A
	return nil
}

// Divides the vector by a scalar
func (v *Vector2D) Div(scalar interface{}) error {
	A, err := getFloat(scalar)
	if err != nil {
		return err
	}
	if A == 0 {
		return fmt.Errorf("divide by zero")
	}
	v.X /= A
	v.Y /= A
	return nil
}

// rotate the vector in the direction of the angle
func (v *Vector2D) Rotate(angle interface{}) error {
	ang, err := getFloat(angle)
	if err != nil {
		return err
	}
	newHeading := v.Heading() + ang
	m := v.Mag()
	v.X = float32(math.Cos(float64(newHeading))) * m
	v.Y = float32(math.Sin(float64(newHeading))) * m
	return nil
}

// Rotate the vector to a specific angle, magnitude remains the same
func (v *Vector2D) SetHeading(angle interface{}) error {
	ang, err := getFloat(angle)
	if err != nil {
		return err
	}
	m := v.Mag()
	v.X = float32(math.Cos(float64(ang))) * m
	v.Y = float32(math.Sin(float64(ang))) * m
	return nil
}

// Calculates the Euclidean distance between two points
// (considering a point as a vector object)
func (v Vector2D) Dist(v2 Vector2D) float32 {
	sV := Sub(v, v2)
	return sV.Mag()
}

// Calculates the dot product with another vector
func (v Vector2D) Dot(v2 Vector2D) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Calculates the cross product with another vector
// ~ give the value of the z axis component
// (in 2D space, the cross product is a vector perpendicular to the two input vectors)
func (v Vector2D) Cross(v2 Vector2D) float32 {
	return v.X*v2.Y - v.Y*v2.X
}

// Calculates and returns the angle with another vector
// Return error if the vectors any vector is zero vector
func (v Vector2D) AngleBetween(v2 Vector2D) (float32, error) {
	m1 := v.Mag()
	m2 := v2.Mag()
	if m1 == 0 || m2 == 0 {
		return 0, fmt.Errorf("cannot calculate angle between zero vectors")
	}
	dotMag := Dot(v, v2) / (m1 * m2)
	angle := math.Acos(math.Min(1, math.Max(-1, float64(dotMag))))
	sign := Cross(v, v2) < 0
	if sign {
		angle = -angle
	}
	return float32(angle), nil
}

// Gets a copy of the vector
func Copy(v Vector2D) Vector2D {
	return Vector2D{v.X, v.Y}
}

// Gives a unit vector in dirction of the vector
func Unit(v Vector2D) Vector2D {
	m := v.Mag()
	if m == 0 {
		return Vector2D{0, 0}
	}
	return Vector2D{v.X / m, v.Y / m}
}

// returns the sum of two vectors
func Add(v1, v2 Vector2D) Vector2D {
	return Vector2D{v1.X + v2.X, v1.Y + v2.Y}
}

// returns the difference of two vectors
func Sub(v1, v2 Vector2D) Vector2D {
	return Vector2D{v1.X - v2.X, v1.Y - v2.Y}
}

// Calculates the dot product of two vectors
func Dot(v1, v2 Vector2D) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Calculates the cross product of two vectors
// ~ give the value of the z axis component
// (in 2D space, the cross product is a vector perpendicular to the two input vectors)
func Cross(v1, v2 Vector2D) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Calculates and returns the angle between two vectors.
// Return error if the vectors any vector is zero vector
func AngleBetween(v1, v2 Vector2D) (float32, error) {
	m1 := v1.Mag()
	m2 := v2.Mag()
	if m1 == 0 || m2 == 0 {
		return 0, fmt.Errorf("cannot calculate angle between zero vectors")
	}
	dotMag := Dot(v1, v2) / (m1 * m2)
	angle := math.Acos(math.Min(1, math.Max(-1, float64(dotMag))))
	sign := Cross(v1, v2) < 0
	if sign {
		angle = -angle
	}
	return float32(angle), nil
}
