package crud

import "time"

func validationStr(obj interface{}) string {
	if obj == nil {
		return ""
	}

	return obj.(string)
}

func validationInt(obj interface{}) int64 {
	if obj == nil {
		return 0
	}

	return obj.(int64)
}

func validationTime(obj interface{}) time.Time {
	if obj == nil {
		return time.Now()
	}

	return obj.(time.Time)
}