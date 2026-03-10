package utils

func RetZeroOnError[T any](data T, err error) (T, error) {
	if err != nil {
		var zero T
		return zero, err
	}
	return data, nil
}
