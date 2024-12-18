package models

type IVFInput struct {
	Age          int
	Weight       int
	Feet         int
	Inches       int
	IVFUsed      string
	Coefficients map[string]interface{}
	ReasonKnown  string
	UseOwnEggs   string
}
