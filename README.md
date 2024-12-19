## Sunfish IVF Success Calculator
This service aims to mimic the functionality of CDC IVF Calculator, which calculates your chance of having a baby using 
In Vitro Fertilization.
This service is designed as a web service with a single endpoint: HTTP GET `/calculate`.  It takes parameters as 
query string and returns a number which corresponds to your chance of having a baby in percents.   

## How to run ## 
run `go run ./cmd/main.go` from the project root.
The endpoint should be available at `http://localhost:8080/calculate`

**Below are some sample requests that you can run to validate the calculator:**
- Using Own Eggs / Did Not Previously Attempt IVF / Known Infertility Reason:
`curl --location 'http://localhost:8080/calculate?age=32&weight=150&feet=5&inches=8&ivf_used=0&gravida=1&tubal_factor=No&male_factor_infertility=No&endometriosis=Yes&ovulatory_disorder=Yes&diminished_ovarian_reserve=No&uterine_factor=No&other_reason=No&unexplained_infertility=No&donotknow=No&eggSource=Own&previous_live_births=1'`
Will return {"success_rate": **62.21** }
- Using Own Eggs / Did Not Previously Attempt IVF / Unknown Infertility Reason:
`curl --location 'http://localhost:8080/calculate?age=32&weight=150&feet=5&inches=8&ivf_used=0&gravida=1&tubal_factor=No&male_factor_infertility=No&endometriosis=No&ovulatory_disorder=No&diminished_ovarian_reserve=No&uterine_factor=No&other_reason=No&unexplained_infertility=No&donotknow=Yes&eggSource=Own&previous_live_births=1'`
  Will return {"success_rate": **59.83** }
- Using Own Eggs / Previously Attempted IVF / Known Infertility Reason:
`curl --location 'http://localhost:8080/calculate?age=32&weight=150&feet=5&inches=8&ivf_used=2&gravida=1&tubal_factor=Yes&male_factor_infertility=No&endometriosis=No&ovulatory_disorder=No&diminished_ovarian_reserve=Yes&uterine_factor=No&other_reason=No&unexplained_infertility=No&donotknow=No&eggSource=Own&previous_live_births=1'`
  Will return {"success_rate": **40.89** }
- Using Donor Eggs / Previously Attempted IVF / Known Infertility Reason:
`curl --location 'http://localhost:8080/calculate?age=32&weight=150&feet=5&inches=8&ivf_used=2&gravida=1&tubal_factor=Yes&male_factor_infertility=No&endometriosis=No&ovulatory_disorder=No&diminished_ovarian_reserve=Yes&uterine_factor=No&other_reason=No&unexplained_infertility=No&donotknow=No&eggSource=Donor&previous_live_births=1'`
  Will return {"success_rate": **51.18** }
- Using Donor Eggs / Previously Attempted IVF / Unknown Infertility Reason:
`curl --location 'http://localhost:8080/calculate?age=32&weight=150&feet=5&inches=8&ivf_used=2&gravida=1&tubal_factor=No&male_factor_infertility=No&endometriosis=No&ovulatory_disorder=No&diminished_ovarian_reserve=No&uterine_factor=No&other_reason=No&unexplained_infertility=No&donotknow=Yes&eggSource=Donor&previous_live_births=1'`
  Will return {"success_rate": **55.8** }

## TODOs ##
- Better test coverage.  The layers are connected via interfaces so it should be easy to mock.  
There is one actual test, however.
- Better input validation.  Most of the validation is there but few edge cases still need to be thought through.  It is 
a good idea to validate in the backend as well as frontend.
- Implement structured logging.  Even tho the logger is initialized and passed down the chain, it is not really used 
and does not allow to easily identify errors and warnings and infos.

## Notes ##
- I tried to mimic the CDC form variables and their values.  Some of them are inconsistent in terms of naming.  For example 
`donotknow` vs `eggSource` vs `previous_live_births`.  I also implemented `Yes/No` (which is case-sensitive) 
and not `true/false` as values for the checkboxes.  In a production env, I would want to make the endpoint less flaky.
- To match the result from the assignment README, I had to round the BMI to a single decimal.  Otherwise, the results were
 very slightly off (by 0.01).


