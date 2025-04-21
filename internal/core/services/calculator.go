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
