package env

import "os"
import "strconv"

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	vaAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return vaAsInt
}

func GetInt64(key string, fallback int64) int64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	vaAsInt64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return vaAsInt64
}
