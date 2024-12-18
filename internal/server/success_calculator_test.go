package server

import (
	"errors"
	"testing"

	"ivf_calculator/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFormulaGetter is a mock implementation of FormulaGetter interface
type MockFormulaGetter struct {
	mock.Mock
}

func (m *MockFormulaGetter) GetFormula(usingOwnEggs string, attemptedIVFPreviously string, isReasonKnown string) (*models.Formula, error) {
	args := m.Called(usingOwnEggs, attemptedIVFPreviously, isReasonKnown)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Formula), args.Error(1)
}

func TestNewSuccessCalculator(t *testing.T) {
	repo := new(MockFormulaGetter)
	calc := NewSuccessCalculator(repo)

	assert.NotNil(t, calc)
	assert.Equal(t, repo, calc.repo)
}

func TestCalculateBMI(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.IVFInput
		expected float64
	}{
		{
			name: "Normal BMI calculation",
			input: &models.IVFInput{
				Weight: 150,
				Feet:   5,
				Inches: 5,
			},
			expected: 25.0,
		},
		{
			name: "Low BMI calculation",
			input: &models.IVFInput{
				Weight: 100,
				Feet:   5,
				Inches: 5,
			},
			expected: 16.6,
		},
	}

	calc := NewSuccessCalculator(nil) // repo not needed for BMI calculation

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateBMI(tt.input)
			assert.InDelta(t, tt.expected, result, 0.1)
		})
	}
}

func TestCalculateSuccess(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.IVFInput
		mockFormula   *models.Formula
		mockError     error
		expected      float64
		expectedError error
	}{
		{
			name: "Basic success calculation with real formula coefficients",
			input: &models.IVFInput{
				UseOwnEggs:  "yes",
				IVFUsed:     "no",
				ReasonKnown: "yes",
				Age:         35,
				Weight:      150,
				Feet:        5,
				Inches:      5,
				Coefficients: map[string]interface{}{
					"tubalFactor":           true,
					"maleFactorInfertility": false,
					"priorPregnancies":      "0",
					"priorLiveBirths":       "0",
				},
			},
			mockFormula: &models.Formula{
				UsingOwnEggs:                "yes",
				AttemptedIVFPreviously:      "no",
				IsReasonForInfertilityKnown: "yes",
				CDCFormula:                  "1-3",
				Coefficients: struct {
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
				}{
					Intercept:      -6.8392144,
					AgeLinear:      0.3347309,
					AgePower:       -0.0003249,
					AgePowerFactor: 2.763313,
					BMILinear:      0.06997997,
					BMIPower:       -0.0015045,
					BMIPowerFactor: 2,
					TubalFactor: map[bool]float64{
						true:  0.09373152,
						false: 0,
					},
					MaleFactorInfertility: map[bool]float64{
						true:  0.24104423,
						false: 0,
					},
					Endometriosis: map[bool]float64{
						true:  0.02773216,
						false: 0,
					},
					OvulatoryDisorder: map[bool]float64{
						true:  0.27949598,
						false: 0,
					},
					DiminishedOvarianReserve: map[bool]float64{
						true:  -0.5780511,
						false: 0,
					},
					UterineFactor: map[bool]float64{
						true:  -0.1354896,
						false: 0,
					},
					OtherReason: map[bool]float64{
						true:  -0.1018557,
						false: 0,
					},
					UnexplainedInfertility: map[bool]float64{
						true:  0.2252616,
						false: 0,
					},
					PriorPregnancies: map[string]float64{
						"0":  0,
						"1":  0.03514055,
						"2+": -0.0059006,
					},
					PriorLiveBirths: map[string]float64{
						"0":  0,
						"1":  0.15787934,
						"2+": 0.03077479,
					},
				},
			},
			expected:      44.39,
			expectedError: nil,
		},
		{
			name: "Error getting formula",
			input: &models.IVFInput{
				UseOwnEggs:  "yes",
				IVFUsed:     "no",
				ReasonKnown: "yes",
			},
			mockFormula:   nil,
			mockError:     errors.New("formula not found"),
			expected:      0,
			expectedError: errors.New("formula not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockFormulaGetter)
			repo.On("GetFormula", tt.input.UseOwnEggs, tt.input.IVFUsed, tt.input.ReasonKnown).
				Return(tt.mockFormula, tt.mockError)

			calc := NewSuccessCalculator(repo)
			result, err := calc.CalculateSuccess(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			repo.AssertExpectations(t)
		})
	}
}
