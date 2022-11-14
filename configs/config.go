package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

type Config struct {
	ApiKey                                     string
	ProjectID                                  uint32
	ListenAddress                              string
	InfuraClientTimeout                        time.Duration
	InfuraClientCacheRefreshInterval           time.Duration
	InfuraClientCacheMaxDurationBeforeEviction time.Duration
}

func NewConfigFromEnv() Config {
	c := Config{
		InfuraClientTimeout:                        getDurationOrDefault("INFURA_CLIENT_TIMEOUT", 10*time.Second),
		InfuraClientCacheRefreshInterval:           getDurationOrDefault("INFURA_CLIENT_CACHE_REFRESH_INTERVAL", time.Minute),
		InfuraClientCacheMaxDurationBeforeEviction: getDurationOrDefault("INFURA_CLIENT_CACHE_MAX_LIFETIME", 10*time.Minute),
	}
	var err error
	var errs []error

	if c.ApiKey, err = getNonEmptyString("INFURA_API_KEY"); err != nil {
		errs = append(errs, err)
	}

	if c.ProjectID, err = getUint[uint32]("INFURA_PROJECT_ID"); err != nil {
		errs = append(errs, err)
	}

	if c.ListenAddress, err = getNonEmptyString("LISTEN_ADDRESS"); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		panic(fmt.Sprintf("conf parsing errors: %+v", errs))
	}

	return c
}

func getUint[T constraints.Unsigned](key string) (T, error) {
	p, err := strconv.ParseUint(os.Getenv(key), 10, int(unsafe.Sizeof(T(0))))
	return T(p), errors.Wrapf(err, "failed to parse %T from [key '%s'/ value '%s']", T(0), key, os.Getenv(key))
}

func getDurationOrDefault(key string, defaultDuration time.Duration) time.Duration {
	if t, err := time.ParseDuration(os.Getenv(key)); err != nil {
		return defaultDuration
	} else {
		return t
	}

}

func getNonEmptyString(key string) (string, error) {
	if v := os.Getenv(key); len(v) == 0 {
		return "", errors.Errorf("env var '%s' is empty", key)
	} else {
		return v, nil
	}
}
