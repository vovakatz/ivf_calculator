package server

import (
	"log"
	"math"

	"ivf_calculator/internal/models"
)

type Config struct {
	Repo   FormulaGetter
	Logger *log.Logger
}

type FormulaGetter interface {
	GetFormula(usingOwnEggs string, attemptedIVFPreviously string, isReasonKnown string) (*models.Formula, error)
}

type SuccessCalculator struct {
	*Config
}

func NewSuccessCalculator(config *Config) *SuccessCalculator {
	return &SuccessCalculator{
		config,
	}
}

// CalculateSuccess calculates the success probability using the formula
func (s *SuccessCalculator) CalculateSuccess(params *models.IVFInput) (float64, error) {
	f, err := s.Repo.GetFormula(params.UseOwnEggs, params.IVFUsed, params.ReasonKnown)
	if err != nil {
		return 0, err
	}

	bmi := s.CalculateBMI(params)
	age := float64(params.Age)

	// Base calculation
	score := f.Coefficients.Intercept +
		f.Coefficients.AgeLinear*age +
		f.Coefficients.AgePower*math.Pow(age, f.Coefficients.AgePowerFactor) +
		f.Coefficients.BMILinear*bmi +
		f.Coefficients.BMIPower*math.Pow(bmi, f.Coefficients.BMIPowerFactor)

	// Add boolean parameters if they exist in the input
	booleanFactors := map[string]*map[bool]float64{
		"tubalFactor":              &f.Coefficients.TubalFactor,
		"maleFactorInfertility":    &f.Coefficients.MaleFactorInfertility,
		"endometriosis":            &f.Coefficients.Endometriosis,
		"ovulatoryDisorder":        &f.Coefficients.OvulatoryDisorder,
		"diminishedOvarianReserve": &f.Coefficients.DiminishedOvarianReserve,
		"uterineFactor":            &f.Coefficients.UterineFactor,
		"otherReason":              &f.Coefficients.OtherReason,
		"unexplainedInfertility":   &f.Coefficients.UnexplainedInfertility,
	}

	for key, coeffMap := range booleanFactors {
		if val, exists := params.Coefficients[key]; exists {
			if boolVal, ok := val.(bool); ok {
				coef := (*coeffMap)[boolVal]
				score += coef
			}
		}
	}

	// Add numeric parameters
	if val, exists := params.Coefficients["priorPregnancies"]; exists {
		if strVal, ok := val.(string); ok {
			score += f.Coefficients.PriorPregnancies[strVal]
		}
	}

	if val, exists := params.Coefficients["priorLiveBirths"]; exists {
		if strVal, ok := val.(string); ok {
			score += f.Coefficients.PriorLiveBirths[strVal]
		}
	}

	// calculate success rate in %
	successRate := 1.0 / (1.0 + math.Exp(-score)) * 100
	// round to 2 decimal digits
	successRate = math.Round(successRate*100) / 100

	return successRate, nil
}

func (s *SuccessCalculator) CalculateBMI(params *models.IVFInput) float64 {
	bmi := float64(params.Weight) / math.Pow(float64(params.Feet*12)+float64(params.Inches), 2) * 703
	// round to single decimal to match assignment results.
	bmi = math.Round(bmi*10) / 10
	return bmi
}
