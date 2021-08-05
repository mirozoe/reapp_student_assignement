package reapp_students_assignement

import "fmt"

const (
	CityValidationError = "City validation error"
	CityDidntFound      = "City wasn't found"
	WorkplaceNotFound   = "Workplace not found in slice"
	UnknownError        = "Unknown error"
)

func ErrorFactory(e, m string) error {
	switch e {
	case CityValidationError:
		return fmt.Errorf("%s %s", CityValidationError, m)
	case CityDidntFound:
		return fmt.Errorf("%s for %s", CityDidntFound, m)
	case WorkplaceNotFound:
		return fmt.Errorf("%s for %s", WorkplaceNotFound, m)
	default:
		return fmt.Errorf("%s", UnknownError)
	}
}
