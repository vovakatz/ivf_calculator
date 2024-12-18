package repo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"ivf_calculator/internal/models"
	"ivf_calculator/internal/utils"
)

type Config struct {
	FilePath string
	Logger   *log.Logger
}

type IVFFormula struct {
	*Config
}

func NewIVFFormula(config *Config) *IVFFormula {
	return &IVFFormula{
		config,
	}
}

// GetFormula reads the CSV file and returns matching formula
func (f *IVFFormula) GetFormula(usingOwnEggs string, attemptedIVFPreviously string, isReasonKnown string) (*models.Formula, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip headers
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading headers: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	for _, record := range records {
		formula := models.Formula{}
		if usingOwnEggs == record[0] && attemptedIVFPreviously == record[1] && isReasonKnown == record[2] {

			// Parse boolean fields
			formula.UsingOwnEggs = usingOwnEggs
			formula.AttemptedIVFPreviously = attemptedIVFPreviously
			formula.IsReasonForInfertilityKnown = isReasonKnown

			formula.CDCFormula = record[3]

			// Parse coefficients
			formula.Coefficients.Intercept = utils.ParseFloat(record[4])
			formula.Coefficients.AgeLinear = utils.ParseFloat(record[5])
			formula.Coefficients.AgePower = utils.ParseFloat(record[6])
			formula.Coefficients.AgePowerFactor = utils.ParseFloat(record[7])
			formula.Coefficients.BMILinear = utils.ParseFloat(record[8])
			formula.Coefficients.BMIPower = utils.ParseFloat(record[9])
			formula.Coefficients.BMIPowerFactor = utils.ParseFloat(record[10])

			// Initialize maps for boolean coefficients
			formula.Coefficients.TubalFactor = map[bool]float64{
				true:  utils.ParseFloat(record[11]),
				false: utils.ParseFloat(record[12]),
			}
			formula.Coefficients.MaleFactorInfertility = map[bool]float64{
				true:  utils.ParseFloat(record[13]),
				false: utils.ParseFloat(record[14]),
			}
			formula.Coefficients.Endometriosis = map[bool]float64{
				true:  utils.ParseFloat(record[15]),
				false: utils.ParseFloat(record[16]),
			}
			formula.Coefficients.OvulatoryDisorder = map[bool]float64{
				true:  utils.ParseFloat(record[17]),
				false: utils.ParseFloat(record[18]),
			}
			formula.Coefficients.DiminishedOvarianReserve = map[bool]float64{
				true:  utils.ParseFloat(record[19]),
				false: utils.ParseFloat(record[20]),
			}
			formula.Coefficients.UterineFactor = map[bool]float64{
				true:  utils.ParseFloat(record[21]),
				false: utils.ParseFloat(record[22]),
			}
			formula.Coefficients.OtherReason = map[bool]float64{
				true:  utils.ParseFloat(record[23]),
				false: utils.ParseFloat(record[24]),
			}
			formula.Coefficients.UnexplainedInfertility = map[bool]float64{
				true:  utils.ParseFloat(record[25]),
				false: utils.ParseFloat(record[26]),
			}

			// Initialize maps for numeric coefficients
			formula.Coefficients.PriorPregnancies = map[string]float64{
				"0":  utils.ParseFloat(record[27]),
				"1":  utils.ParseFloat(record[28]),
				"2+": utils.ParseFloat(record[29]),
			}
			formula.Coefficients.PriorLiveBirths = map[string]float64{
				"0":  utils.ParseFloat(record[30]),
				"1":  utils.ParseFloat(record[31]),
				"2+": utils.ParseFloat(record[32]),
			}

			return &formula, nil
		}
	}

	return nil, fmt.Errorf("no matching formula found for the given parameters")
}
