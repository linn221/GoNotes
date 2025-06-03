package formscanner

import "net/http"

// do both scanning and validation
func Scan[T any](r *http.Request,
	inputName string,
	ptr *T,
	parseFunc func(*http.Request, string) (T, bool, error),
	validateFuncs ...ValidateFunc[T],
) error {
	v, omitted, err := parseFunc(r, inputName)
	if err != nil {
		return err
	}
	*ptr = v
	if !omitted {
		if len(validateFuncs) > 0 {
			for _, f := range validateFuncs {
				if err := f(v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
