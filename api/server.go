package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"ivf_calculator/internal/models"
	"ivf_calculator/internal/utils"
)

type Config struct {
	Port       string
	Logger     *log.Logger
	IVFService IVFCalculator
}

type Server struct {
	*Config
}

type IVFCalculator interface {
	CalculateSuccess(params *models.IVFInput) (float64, error)
}

func New(config *Config) *Server {
	return &Server{
		config,
	}
}

func (s *Server) Start() {
	// Register handlers
	http.HandleFunc("/calculate", s.CalculateIVFSuccessHandler)

	// Start server
	s.Logger.Printf("Starting server on port %s", s.Port)
	if err := http.ListenAndServe(s.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) CalculateIVFSuccessHandler(w http.ResponseWriter, r *http.Request) {
	s.Logger.Printf("Received request from %s", r.RemoteAddr)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query()
	input, err := s.validateInput(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, err := s.IVFService.CalculateSuccess(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		SuccessRate float64 `json:"success_rate"`
	}{
		SuccessRate: rate,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) validateInput(params url.Values) (*models.IVFInput, error) {
	input := &models.IVFInput{}

	if ageStr := params.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return nil, err
		}
		if age < 20 || age > 50 {
			return nil, fmt.Errorf("age must be between 20 and 50. Got %s", ageStr)
		}
		input.Age = age
	}

	if weightStr := params.Get("weight"); weightStr != "" {
		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			return nil, err
		}
		if weight < 80 || weight > 300 {
			return nil, fmt.Errorf("weight must be between 80 and 300. Got %s", weightStr)
		}
		input.Weight = weight
	}

	if feetStr := params.Get("feet"); feetStr != "" {
		feet, err := strconv.Atoi(feetStr)
		if err != nil {
			return nil, err
		}
		input.Feet = feet
	}

	if inchesStr := params.Get("inches"); inchesStr != "" {
		inches, err := strconv.Atoi(inchesStr)
		if err != nil {
			return nil, err
		}
		input.Inches = inches
	}

	// Populate coefficients map
	input.Coefficients = make(map[string]interface{})
	priorPregnanciesStr := params.Get("gravida")
	if priorPregnanciesStr != "" {
		valGravida := []string{"0", "1", "2+"}
		if utils.Contains(valGravida, priorPregnanciesStr) {
			input.Coefficients["priorPregnancies"] = priorPregnanciesStr
		} else {
			return nil, fmt.Errorf("gravida has invalid value %s", priorPregnanciesStr)
		}
	} else {
		return nil, fmt.Errorf("gravida is required")
	}

	priorLiveBirthsStr := params.Get("previous_live_births")
	if priorLiveBirthsStr != "" {
		valLiveBirths := []string{"0", "1", "2+"}
		if utils.Contains(valLiveBirths, priorLiveBirthsStr) {
			if priorLiveBirthsStr > priorPregnanciesStr {
				return nil, fmt.Errorf("previous_live_births can't be greater then gravida")
			}
			input.Coefficients["priorLiveBirths"] = priorLiveBirthsStr
		} else {
			return nil, fmt.Errorf("previous_live_births has invalid value %s", priorLiveBirthsStr)
		}
	} else {
		return nil, fmt.Errorf("previous_live_births is required")
	}

	factors := map[string]string{
		"tubal_factor":               "tubalFactor",
		"male_factor_infertility":    "maleFactorInfertility",
		"endometriosis":              "endometriosis",
		"ovulatory_disorder":         "ovulatoryDisorder",
		"diminished_ovarian_reserve": "diminishedOvarianReserve",
		"uterine_factor":             "uterineFactor",
		"other_reason":               "otherReason",
	}

	known_reasons := false
	for paramName, coefficientName := range factors {
		if value, err := processKnownReasons(&params, paramName); err != nil {
			return nil, err
		} else {
			input.Coefficients[coefficientName] = value
			if value {
				known_reasons = true
			}
		}
	}

	unexplainedInfertilitySel := false
	if noReasonStr := params.Get("unexplained_infertility"); noReasonStr != "" {
		switch noReasonStr {
		case `Yes`:
			input.Coefficients["unexplainedInfertility"] = true
			unexplainedInfertilitySel = true
		case `No`:
			input.Coefficients["unexplainedInfertility"] = false
		default:
			return nil, fmt.Errorf("unexplained_infertility has invalid value %s", noReasonStr)
		}
	}

	noReasonSel := false
	if noReasonStr := params.Get("donotknow"); noReasonStr != "" {
		switch noReasonStr {
		case `Yes`:
			input.ReasonKnown = "FALSE"
			noReasonSel = true
		case `No`:
			input.ReasonKnown = "TRUE"
		default:
			return nil, fmt.Errorf("no_reason has invalid value %s", noReasonStr)
		}
	}

	if !utils.OnlyOneTrue(known_reasons, unexplainedInfertilitySel, noReasonSel) {
		return nil, fmt.Errorf("known_reasons OR unexplained_infertility OR no_reason is required")
	}

	if useOwnEggsStr := params.Get("eggSource"); useOwnEggsStr != "" {
		switch useOwnEggsStr {
		case `Own`:
			input.UseOwnEggs = "TRUE"
		case `Donor`:
			input.UseOwnEggs = "FALSE"
		default:
			return nil, fmt.Errorf("eggSource has invalid value %s", useOwnEggsStr)
		}
	}

	if ivfusedStr := params.Get("ivf_used"); ivfusedStr != "" {
		if input.UseOwnEggs == "FALSE" {
			input.IVFUsed = "N/A"
		} else {
			valIVFUsed := []string{"0", "1", "2", "3+"}
			if utils.Contains(valIVFUsed, ivfusedStr) {
				if ivfusedStr == "0" {
					input.IVFUsed = "FALSE"
				} else {
					input.IVFUsed = "TRUE"
				}
			} else {
				return nil, fmt.Errorf("ivf_used has invalid value %s", ivfusedStr)
			}
		}
	} else {
		return nil, fmt.Errorf("ivf_used is required")
	}

	return input, nil
}

func processKnownReasons(params *url.Values, paramName string) (bool, error) {
	if paramStr := params.Get(paramName); paramStr != "" {
		switch paramStr {
		case `Yes`:
			return true, nil
		case `No`:
			return false, nil
		default:
			return false, fmt.Errorf("%s has invalid value %s", paramName, paramStr)
		}
	}
	return false, fmt.Errorf("%s is required", paramName)
}
