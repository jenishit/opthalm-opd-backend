package dto

// None of these numeric fields use binding:"required" — 0 is a physically
// valid value for every one of them (a plano sphere, an axis of exactly 0
// degrees, vertexing to the corneal plane at 0mm, etc.), and Go's validator
// treats a zero value as "not provided" for "required", which would wrongly
// reject those legitimate inputs.

type TranspositionReq struct {
	Sphere   float64 `json:"sphere"`
	Cylinder float64 `json:"cylinder"`
	Axis     int     `json:"axis" binding:"min=0,max=180"`
}

type TranspositionRes struct {
	Sphere   float64 `json:"sphere"`
	Cylinder float64 `json:"cylinder"`
	Axis     int     `json:"axis"`
}

type SphericalEquivalentReq struct {
	Sphere   float64 `json:"sphere"`
	Cylinder float64 `json:"cylinder"`
}

type NearAddReq struct {
	DistanceSphere float64 `json:"distance_sphere"`
	AddPower       float64 `json:"add_power"`
}

type VertexDistanceReq struct {
	Power          float64 `json:"power"`
	FromDistanceMM float64 `json:"from_distance_mm"`
	ToDistanceMM   float64 `json:"to_distance_mm"`
}

type TelescopeFOVReq struct {
	TrueFOVDegrees float64 `json:"true_fov_degrees"`
	Magnification  float64 `json:"magnification"`
}

type ScalarResultRes struct {
	Result float64 `json:"result"`
}
