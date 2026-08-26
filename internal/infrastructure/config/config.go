// Package implementation for policy-driven content moderation and human review.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var ErrInvalidConfig = errors.New("invalid config")

type Config struct {
	HTTPAddr        string
	ShutdownSeconds int
	Environment     string
}

func Load() Config {
	c := Config{HTTPAddr: ":8083", ShutdownSeconds: 10, Environment: "local"}
	path := os.Getenv("CONTENT_MODERATION_CONFIG")
	if path == "" {
		path = "configs/config.yaml"
	}
	if parsed, err := ParseYAML(path); err == nil {
		c = Merge(c, parsed)
	}
	c = ApplyEnvironment(c)
	if v := os.Getenv("CONTENT_MODERATION_SHUTDOWN_SECONDS"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.ShutdownSeconds = n
		}
	}
	return c
}

func LoadChecked() (Config, error) {
	path := os.Getenv("CONTENT_MODERATION_CONFIG")
	if path == "" {
		path = "configs/config.yaml"
	}
	if _, err := ParseYAML(path); err != nil {
		return Config{}, err
	}
	return Load(), nil
}
func ParseTags(s string) map[string]string {
	out := map[string]string{}
	for _, p := range strings.Split(s, ",") {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out
}
