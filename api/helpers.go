package api


//
//  Errors
//
type FormattedError struct {
  Error string `json:"error"`
}


func FormatError(err error) FormattedError {
  return FormattedError{ Error: err.Error() }
}
