package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	props   map[string]string
	propsMu sync.RWMutex
)

func LoadConfig(path string) error {
	propsMu.Lock()
	defer propsMu.Unlock()

	props = make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if envVal, ok := os.LookupEnv(strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(key, ".", "_"), "-", "_"))); ok {
			value = envVal
		}

		props[key] = value
	}
	return scanner.Err()
}

func GetString(key string, defaultVal string) string {
	propsMu.RLock()
	defer propsMu.RUnlock()

	if v, ok := props[key]; ok && v != "" && v != "FROM_VAULT" {
		return v
	}
	return defaultVal
}

func GetInt(key string, defaultVal int) int {
	s := GetString(key, "")
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

func GetBool(key string, defaultVal bool) bool {
	s := GetString(key, "")
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return defaultVal
	}
	return v
}
