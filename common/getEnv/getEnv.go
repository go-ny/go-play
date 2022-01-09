package getEnv

import "os"

//EnvWithKey : get env value
func EnvWithKey(key string) string {
	return os.Getenv(key)
}
