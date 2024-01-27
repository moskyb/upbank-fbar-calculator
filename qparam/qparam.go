package qparam

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/oleiade/reflections"
)

const tag = "qparam"

func Merge(vals ...url.Values) url.Values {
	out := make(url.Values)

	for _, val := range vals {
		for k, v := range val {
			out[k] = v
		}
	}

	return out
}

// Marshal takes a struct annotated with `qparam:"..."` tags and returns a url.Values
// empty values are ignored
// multiple values are ignored, only the final value will be set
func Marshal(v any) (url.Values, error) {
	fields, err := reflections.Fields(v)
	if err != nil {
		return nil, fmt.Errorf("gettings fields for type %T: %w", v, err)
	}

	vals := make(url.Values, len(fields))

	for _, field := range fields {
		tagVal, err := reflections.GetFieldTag(v, field, tag)
		if err != nil {
			return nil, fmt.Errorf("getting tag %s for field %s of type %T: %w", tag, field, v, err)
		}

		if tagVal == "-" {
			continue
		}

		val, err := reflections.GetField(v, field)
		if err != nil {
			return nil, fmt.Errorf("getting field %s for type %T: %w", field, v, err)
		}

		qVal, ok := encodeQueryVal(val)
		if ok {
			vals.Set(tagVal, qVal)
		}
	}

	return vals, nil
}

func encodeQueryVal(v any) (val string, use bool) {
	switch v := v.(type) {
	case string:
		if v == "" {
			return "", false
		}

		return v, true

	case int:
		if v == 0 {
			return "", false
		}

		return strconv.Itoa(v), true

	case bool:
		return strconv.FormatBool(v), true

	case time.Time:
		if v.IsZero() {
			return "", false
		}

		return v.Format(time.RFC3339), true

	default:
		return fmt.Sprintf("%v", v), true
	}
}
