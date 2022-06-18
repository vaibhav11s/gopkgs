// Package vector provides a simple 3D vector class.
package vector

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Vector struct {
	X, Y, Z float32
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Creates a new 3D vector.
// Three dimensional Euclidean vector.
func New(x, y, z float32) *Vector {
	return &Vector{x, y, z}
}

func x() *Vector {
	return &Vector{1, 0, 0}
}

func y() *Vector {
	return &Vector{0, 1, 0}
}

func z() *Vector {
	return &Vector{0, 0, 1}
}

func zero() *Vector {
	return &Vector{0, 0, 0}
}

// Make a new 3D vector from a pair of azimuth and zenith angles.
// https://en.wikipedia.org/wiki/Spherical_coordinate_system
func FromAngles(thetha, phi float32, length ...float32) *Vector {
	var l float64 = 1
	if len(length) >= 1 {
		l = float64(length[0])
	}
	cosPhi := math.Cos(float64(phi))
	sinPhi := math.Sin(float64(phi))
	cosTheta := math.Cos(float64(thetha))
	sinTheta := math.Sin(float64(thetha))
	return &Vector{
		X: float32(l * cosTheta * sinPhi),
		Y: float32(l * sinTheta * sinPhi),
		Z: float32(l * cosPhi),
	}
}

// Makes a random 3D vector of given lenght (default 1)
func Random(length ...float32) *Vector {
	var l float32 = 1
	if len(length) >= 1 {
		l = length[0]
	}
	thetha := rand.Float32() * 2 * math.Pi
	phi := rand.Float32() * 2 * math.Pi
	return FromAngles(thetha, phi, l)
}

// String representation of vector
func (v *Vector) String() string {
	return fmt.Sprintf("{X: %v, Y: %v, Z: %v}", v.X, v.Y, v.Z)
}

// Checks whether two vectors are equal.
// optional tolerence value can be passed as a parameter to check for equality
// within a tolerance.
// abs(v.x - v2.x) < tolerance && abs(v.y - v2.y) < tolerance && abs(v.z - v2.z) < tolerance
func (v *Vector) Equal(v2 *Vector, tolerance ...float32) bool {
	var t float32 = 1e-7
	if len(tolerance) >= 1 {
		t += tolerance[0]
	}
	if diff := float32(math.Abs(float64(v.X - v2.X))); diff > t {
		return false
	}
	if math.Abs(float64(v.Y-v2.Y)) > float64(t) {
		return false
	}
	if math.Abs(float64(v.Z-v2.Z)) > float64(t) {
		return false
	}
	return true
}

func isZero(v *Vector) bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Gets a copy of the vector
func (v *Vector) Copy() *Vector {
	return &Vector{v.X, v.Y, v.Z}
}

// Gets a copy of the vector
func Copy(v *Vector) *Vector {
	return &Vector{v.X, v.Y, v.Z}
}

// Assigns the values of given vector to the vector.
// Similar to copy, but no new vector is create
func (v1 *Vector) Assign(v2 *Vector) *Vector {
	v1.X = v2.X
	v1.Y = v2.Y
	v1.Z = v2.Z
	return v1
}

// Calculates the magnitude (length) of the vector and returns the result as a float
// this is simply the equation sqrt(x*x + y*y + z*z)
func (v *Vector) Mag() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Calculates the squared magnitude of the vector and returns the result as a float
// this is simply the equation (x*x + y*y + z*z)
func (v *Vector) MagSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Normalize the vector to length 1 (make it a unit vector).
// Modify + Returns self
func (v *Vector) Normalize() *Vector {
	mag := v.Mag()
	if mag != 0 {
		v.X /= mag
		v.Y /= mag
		v.Z /= mag
	}
	return v
}

// Gives a unit vector in dirction of the vector
func Unit(v *Vector) *Vector {
	m := v.Mag()
	if m == 0 {
		return &Vector{0, 0, 0}
	}
	return &Vector{v.X / m, v.Y / m, v.Z / m}
}

// Set the magnitude of the vector to the given value.
// Modify + Returns self
func (v *Vector) Resize(mag float32) *Vector {
	v.Normalize()
	v.Mult(mag)
	return v
}

// add a vector to the current vector.
// Modify + Returns self
func (v *Vector) Add(v2 *Vector) *Vector {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

// returns the sum of two vectors
func Add(v1, v2 *Vector) *Vector {
	return &Vector{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

// subtract a vector from the current vector.
// Modify + Returns self
func (v *Vector) Sub(v2 *Vector) *Vector {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

// returns the difference of two vectors
func Sub(v1, v2 *Vector) *Vector {
	return &Vector{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

// Multiplies the vector by a scalar.
// Modify + Returns self
func (v *Vector) Mult(scalar float32) *Vector {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

// Calculates the Euclidean distance between two points
// (considering a point as a vector object)
func (v *Vector) Dist(v2 *Vector) float32 {
	return Dist(v, v2)
}

// Calculates the Euclidean distance between two points
// (considering a point as a vector object)
func Dist(v1, v2 *Vector) float32 {
	return Sub(v1, v2).Mag()
}

// Calculates the dot product with another vector
func (v *Vector) Dot(v2 *Vector) float32 {
	return Dot(v, v2)
}

// Calculates the dot product of two vectors
func Dot(v1, v2 *Vector) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Calculates the cross product with another vector
func (v *Vector) Cross(v2 *Vector) *Vector {
	return Cross(v, v2)
}

// Calculates the cross product of two vectors
func Cross(v1, v2 *Vector) *Vector {
	return &Vector{v1.Y*v2.Z - v1.Z*v2.Y, v1.Z*v2.X - v1.X*v2.Z, v1.X*v2.Y - v1.Y*v2.X}
}

// Calculates and returns the angle with another vector
// Returns NaN if any vector is a zero vector
func (v *Vector) Angle(v2 *Vector) float32 {
	return Angle(v, v2)
}

// Calculates and returns the angle between two vectors.
// Returns NaN if any vector is a zero vector
func Angle(v1, v2 *Vector) float32 {
	m1 := v1.Mag()
	m2 := v2.Mag()
	if m1 == 0 || m2 == 0 {
		return float32(math.NaN())
	}
	return float32(math.Acos(float64(Dot(v1, v2) / (m1 * m2))))
}

// Calculate the azimuth and zenith angles.
// https://en.wikipedia.org/wiki/Spherical_coordinate_system
func (v *Vector) Heading() (theta, phi float32) {
	m := v.Mag()
	theta = float32(math.Atan2(float64(v.Y), float64(v.X)))
	if m == 0 {
		phi = float32(math.NaN())
		return
	}
	phi = float32(math.Acos(float64(v.Z / m)))
	return
}

// Rotate the vector to a specific angle. magnitude remains the same.
// Modify + Returns self
// https://en.wikipedia.org/wiki/Spherical_coordinate_system
func (v *Vector) SetHeading(thetha, phi float32) *Vector {
	l := float64(v.Mag())
	cosPhi := math.Cos(float64(phi))
	sinPhi := math.Sin(float64(phi))
	cosTheta := math.Cos(float64(thetha))
	sinTheta := math.Sin(float64(thetha))
	v.X = float32(l * cosTheta * sinPhi)
	v.Y = float32(l * sinTheta * sinPhi)
	v.Z = float32(l * cosPhi)
	return v
}

func rotateOnPlane(v, normal *Vector, angle float32) *Vector {
	// v dot n = 0
	sin := float32(math.Sin(float64(angle)))
	cos := float32(math.Cos(float64(angle)))
	nv := Cross(Unit(normal), v)
	nv.Mult(sin)
	V := v.Copy().Mult(cos)
	V.Add(nv)
	return V
}

func (v *Vector) rotateOnPlane(normal *Vector, angle float32) *Vector {
	v.Assign(rotateOnPlane(v, normal, angle))
	return v
}

// Give the component of the given vector parallel and perpendicular to the axis
func (v *Vector) Component(axis *Vector) (parallel, perpendicular *Vector) {
	if isZero(axis) {
		return zero(), zero()
	}
	parallel = axis.Copy().Normalize()
	parallel.Mult(Dot(v, parallel))
	perpendicular = Sub(v, parallel)
	return
}

// Rotates the given vector around the axis by given angle
func RotateAlongAxis(v, axis *Vector, angle float32) *Vector {
	if isZero(axis) {
		return v
	}
	parallel, perpendicular := v.Component(axis)
	perpendicular.rotateOnPlane(axis, angle)
	parallel.Add(perpendicular)
	return parallel
}

// Rotates the given vector around the axis by given angle
// https://math.stackexchange.com/questions/511370/how-to-rotate-one-vector-about-another
func (v *Vector) RotateAlongAxis(axis *Vector, angle float32) *Vector {
	v.Assign(RotateAlongAxis(v, axis, angle))
	return v
}

// Gives the reflection of vector from the given plane(normal vector)
func ReflectThroughPlane(v, normal *Vector) *Vector {
	if isZero(normal) {
		return v
	}
	n := normal.Copy().Normalize()
	return Sub(v, n.Mult(2*Dot(v, n)))
}

// Gives the reflection of vector from the given plane(normal vector)
func (v *Vector) ReflectThroughPlane(normal *Vector) *Vector {
	v.Assign(ReflectThroughPlane(v, normal))
	return v
}

func lerpf(a, b, t float32) float32 {
	return a + (b-a)*t
}

// Linear interpolate the vector to another vector
func Lerp(v1, v2 *Vector, t float32) *Vector {
	return &Vector{lerpf(v1.X, v2.X, t), lerpf(v1.Y, v2.Y, t), lerpf(v1.Z, v2.Z, t)}
}

// Linear interpolate the vector to another vector. i/n = t
func Lerp2(v1, v2 *Vector, n, i int) *Vector {
	r := float32(i) / float32(n)
	return Lerp(v1, v2, r)
}
