package env

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func LoadEnvironment(baseDir string, files ...string) {
	env := os.Getenv("ENV")
	if env == `` || env == `dev` {
		env = "development"
		os.Setenv(`ENV`, env)
	}

	files = append([]string{`.env.` + env + `.local`, `.env.` + env, `.env.local`, `.env`}, files...)

	for _, file := range files {
		if err := godotenv.Load(baseDir + file); err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(`Error loading environment file(s):` + err.Error())
		}
	}
}

func IsEnv(s string) bool {
	return GetString(`ENV`, `development`) == s
}

func GetString(key, def string) string {
	v := os.Getenv(key)
	if v == `` {
		return def
	}
	return v
}

func GetInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

func GetFloat(key string, defaultValue float64) float64 {
	v := os.Getenv(key)
	if v == `` {
		return defaultValue
	}

	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

func GetUInt(key string, defaultValue uint64) uint64 {
	return uint64(GetInt(key, int(defaultValue)))
}
