package services

// CalculatorService groups the stateless optical math helpers used by
// front-desk/optometry staff. There's nothing to persist or look up, so
// unlike every other service here it has no repository dependency.
type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

// TransposeCylinder converts a prescription between plus-cylinder and
// minus-cylinder notation: new sphere = sphere + cylinder, new cylinder is
// negated, new axis is rotated 90 degrees (wrapped into [0, 180)).
func (c *CalculatorService) TransposeCylinder(sphere, cylinder float64, axis int) (newSphere, newCylinder float64, newAxis int) {
	newSphere = sphere + cylinder
	newCylinder = -cylinder
	newAxis = axis + 90
	if newAxis >= 180 {
		newAxis -= 180
	}
	return
}

// SphericalEquivalent = sphere + (cylinder / 2).
func (c *CalculatorService) SphericalEquivalent(sphere, cylinder float64) float64 {
	return sphere + cylinder/2
}

// NearAdd computes the near prescription for one eye: distance sphere + add power.
func (c *CalculatorService) NearAdd(distanceSphere, addPower float64) float64 {
	return distanceSphere + addPower
}

// VertexDistanceAdjust converts a lens power measured at one vertex distance
// to its equivalent power at another (distances in millimeters, powers in
// diopters), using the standard F' = F / (1 - d*F) formula.
func (c *CalculatorService) VertexDistanceAdjust(power float64, fromDistanceMM, toDistanceMM float64) float64 {
	deltaMeters := (toDistanceMM - fromDistanceMM) / 1000
	denominator := 1 - deltaMeters*power
	if denominator == 0 {
		return power
	}
	return power / denominator
}

// TelescopeFOV computes the apparent field of view (degrees) of a telescope
// given its true (objective-side) field of view and magnification.
func (c *CalculatorService) TelescopeFOV(trueFOVDegrees, magnification float64) float64 {
	return trueFOVDegrees * magnification
}
