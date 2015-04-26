package api

type FormattedError struct {
	Error string `json:"error"`
}

func FormatError(err error) FormattedError {
	return FormattedError{Error: err.Error()}
}

func IsNoResultsError(err string) bool {
	var msg string = "sql: no rows in result set"

	if err == msg {
		return true
	} else {
		return false
	}
}
