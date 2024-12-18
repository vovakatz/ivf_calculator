package models

// Formula represents a single IVF success formula with all its parameters
type Formula struct {
	UsingOwnEggs                string
	AttemptedIVFPreviously      string
	IsReasonForInfertilityKnown string
	CDCFormula                  string
	Coefficients                struct {
		Intercept                float64
		AgeLinear                float64
		AgePower                 float64
		AgePowerFactor           float64
		BMILinear                float64
		BMIPower                 float64
		BMIPowerFactor           float64
		TubalFactor              map[bool]float64
		MaleFactorInfertility    map[bool]float64
		Endometriosis            map[bool]float64
		OvulatoryDisorder        map[bool]float64
		DiminishedOvarianReserve map[bool]float64
		UterineFactor            map[bool]float64
		OtherReason              map[bool]float64
		UnexplainedInfertility   map[bool]float64
		PriorPregnancies         map[string]float64
		PriorLiveBirths          map[string]float64
	}
}
