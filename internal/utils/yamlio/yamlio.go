// Package yamlio provides cross-platform YAML I/O with guaranteed UTF-8 encoding.
// Mirrors src/apm_cli/utils/yaml_io.py.
//
// NOTE: Full YAML parsing requires an external library (gopkg.in/yaml.v3). This
// package provides the API surface and a minimal implementation that handles the
// common cases APM uses (string/int/bool values, no anchors/aliases). Production
// callers that need full YAML support should build with gopkg.in/yaml.v3 and swap
// the internal parseYAML / marshalYAML implementations.
package yamlio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadYAML reads a YAML file and returns the parsed data as a flat map.
// Returns nil for empty files. Returns an error on failure.
func LoadYAML(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	return parseSimpleYAML(string(data))
}

// DumpYAML writes data to a YAML file with UTF-8 encoding.
func DumpYAML(data any, path string) error {
	out, err := YAMLToStr(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(out), 0o644)
}

// YAMLToStr serializes a map[string]any to a minimal YAML string.
func YAMLToStr(data any) (string, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return fmt.Sprintf("%v\n", data), nil
	}
	var sb strings.Builder
	for k, v := range m {
		sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}
	return sb.String(), nil
}

// parseSimpleYAML handles flat "key: value" YAML (no nesting, anchors, or sequences).
func parseSimpleYAML(content string) (map[string]any, error) {
	result := map[string]any{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		result[key] = val
	}
	return result, scanner.Err()
}
