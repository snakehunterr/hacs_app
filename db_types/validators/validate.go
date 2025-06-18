package validators

import (
	"strconv"
	"time"

	api_errors "github.com/snakehunterr/hacs_app/db_types/errors"
)

const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

func Int64(name, value string, emptyCheck bool) (i int64, err *api_errors.APIError) {
	if emptyCheck && len(value) == 0 {
		return i, api_errors.NewErrEmptyParam(name)
	}
	v, e := strconv.ParseInt(value, 10, 64)
	if e != nil {
		return i, api_errors.NewErrIncorrectParam(name)
	}
	return v, nil
}

func Uint8(name, value string, emptyCheck bool) (u uint8, err *api_errors.APIError) {
	if emptyCheck && len(value) == 0 {
		return u, api_errors.NewErrEmptyParam(name)
	}
	v, e := strconv.ParseUint(value, 10, 8)
	if e != nil {
		return u, api_errors.NewErrIncorrectParam(name)
	}
	return uint8(v), nil
}

func Float64(name, value string, emptyCheck bool) (f float64, err *api_errors.APIError) {
	if emptyCheck && len(value) == 0 {
		return f, api_errors.NewErrEmptyParam(name)
	}
	v, e := strconv.ParseFloat(value, 64)
	if e != nil {
		return f, api_errors.NewErrIncorrectParam(name)
	}
	return v, nil
}

func Date(name, value string, emptyCheck bool) (d time.Time, err *api_errors.APIError) {
	if emptyCheck && len(value) == 0 {
		return time.Time{}, api_errors.NewErrEmptyParam(name)
	}
	v, e := time.Parse(DATE_FORMAT, value)
	if e != nil {
		return v, api_errors.NewErrIncorrectParam(name)
	}
	return v, nil
}

func Bool(name, value string, emptyCheck bool) (b bool, err *api_errors.APIError) {
	if emptyCheck && len(value) == 0 {
		return b, api_errors.NewErrEmptyParam(name)
	}
	v, e := strconv.ParseBool(value)
	if e != nil {
		return b, api_errors.NewErrIncorrectParam(name)
	}
	return v, nil
}
