package drift

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/driftwatch/internal/docker"
)

// DriftResult holds the comparison outcome for a single service.
type DriftResult struct {
	Service  string
	Drifted  bool
	Reasons  []string
}

// Detect compares running containers against the declared compose services.
// specServices maps service name -> declared env vars (key=value).
// containers is the list of currently running containers.
func Detect(specServices map[string][]string, containers []docker.ContainerInfo) []DriftResult {
	results := make([]DriftResult, 0, len(specServices))

	runningByService := indexByService(containers)

	for service, declaredEnv := range specServices {
		result := DriftResult{Service: service}

		info, found := runningByService[service]
		if !found {
			result.Drifted = true
			result.Reasons = append(result.Reasons, "service not running")
			results = append(results, result)
			continue
		}

		if reasons := compareEnv(declaredEnv, info.Env); len(reasons) > 0 {
			result.Drifted = true
			result.Reasons = append(result.Reasons, reasons...)
		}

		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Service < results[j].Service
	})
	return results
}

// indexByService builds a map from compose service label to ContainerInfo.
func indexByService(containers []docker.ContainerInfo) map[string]docker.ContainerInfo {
	m := make(map[string]docker.ContainerInfo, len(containers))
	for _, c := range containers {
		if svc, ok := c.Labels["com.docker.compose.service"]; ok {
			m[svc] = c
		}
	}
	return m
}

// compareEnv returns drift reasons for environment variable mismatches.
func compareEnv(declared []string, running map[string]string) []string {
	var reasons []string
	for _, entry := range declared {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, wantVal := parts[0], parts[1]
		gotVal, exists := running[key]
		if !exists {
			reasons = append(reasons, fmt.Sprintf("env %q missing in running container", key))
		} else if gotVal != wantVal {
			reasons = append(reasons, fmt.Sprintf("env %q: want %q, got %q", key, wantVal, gotVal))
		}
	}
	return reasons
}
