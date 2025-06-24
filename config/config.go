package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv() (map[string]string, error) {
	env := make(map[string]string)

	file, err := os.Open("config/.env")
	if err != nil {
		return nil, fmt.Errorf("Не удалось открыть .env файл: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "=")
		env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return env, nil
}
