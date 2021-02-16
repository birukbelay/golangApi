# data 

``type Input struct {
	Values  url.Values
	VErrors ValidationErrors
	CSRF    string
}
``
- min length 
- required
- matches pattern
- Password matches function


# errors 

type ValidationErrors map[string][]string
- Add(field, message string){}
- Get(field string) string {}

# form_validator

- FormMovieValidator(w http.ResponseWriter, values url.Values)
- FormGenresValidator(values url.Values)