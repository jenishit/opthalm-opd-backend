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

// IOLPowerSRKII estimates the intraocular lens power (diopters) to implant
// for a target post-op refraction, using the SRK-II regression formula —
// the standard quick-reference formula for cataract/IOL planning, refining
// the original SRK formula with an axial-length-dependent correction so it
// stays accurate for short and long eyes, not just average ones.
//
//	axialLengthMM: biometry axial length (mm), typically 20-26mm
//	avgKDiopters: average corneal power from keratometry (D)
//	aConstant: the lens manufacturer's A-constant for the IOL model
//	targetRefractionD: desired post-op spherical equivalent (D); 0 for emmetropia
func (c *CalculatorService) IOLPowerSRKII(axialLengthMM, avgKDiopters, aConstant, targetRefractionD float64) float64 {
	correctedAConstant := aConstant
	switch {
	case axialLengthMM < 20:
		correctedAConstant += 3
	case axialLengthMM < 21:
		correctedAConstant += 2
	case axialLengthMM < 22:
		correctedAConstant += 1
	case axialLengthMM >= 24.5:
		correctedAConstant -= 0.5
	}

	emmetropicPower := correctedAConstant - 0.9*avgKDiopters - 2.5*axialLengthMM
	if targetRefractionD == 0 {
		return emmetropicPower
	}

	// Adjust for a non-zero target refraction: each diopter of desired
	// myopia/hyperopia shifts the required IOL power by roughly 1.0-1.25x
	// itself depending on axial length; SRK-II uses a fixed correction
	// factor of 1.0 for the emmetropic formula's refraction adjustment.
	return emmetropicPower - targetRefractionD*1.0
}
