// Package targetscmd implements the "apm targets" command, which inspects
// and displays the resolved target list for the current project.
package targetscmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/githubnext/apm/internal/core/targetdetection"
)

// TargetRow represents a single target in the output table.
type TargetRow struct {
	Target    string `json:"target"`
	Status    string `json:"status"`
	Source    string `json:"source,omitempty"`
	DeployDir string `json:"deploy_dir"`
	Needs     string `json:"needs,omitempty"`
}

// Run implements the "apm targets" command.
// asJSON prints machine-readable JSON; showAll includes meta-targets.
func Run(asJSON, showAll bool) error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("targetscmd: getwd: %w", err)
	}

	resolved, err := targetdetection.ResolveTargets(projectRoot, nil, nil)
	active := map[string]bool{}
	if err != nil {
		// Fall through with empty active set
	} else {
		for _, t := range resolved.Targets {
			active[t] = true
		}
	}

	signals := targetdetection.DetectSignals(projectRoot)
	signalSources := map[string]string{}
	for _, s := range signals {
		signalSources[s.Target] = s.Source
	}

	rows := make([]TargetRow, 0, len(targetdetection.CanonicalTargetsOrdered))
	for _, name := range targetdetection.CanonicalTargetsOrdered {
		row := TargetRow{
			Target:    name,
			DeployDir: targetdetection.CanonicalDeployDirs[name],
		}
		if active[name] {
			row.Status = "active"
			row.Source = signalSources[name]
		} else {
			row.Status = "inactive"
			row.Needs = targetdetection.CanonicalSignal[name]
		}
		rows = append(rows, row)
	}

	if asJSON {
		if showAll {
			metaStatus := "inactive"
			if active["agent-skills"] {
				metaStatus = "active"
			}
			rows = append(rows, TargetRow{
				Target:    "agent-skills",
				Status:    metaStatus,
				DeployDir: ".agents/",
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}

	fmt.Printf("  %-12s  %-10s  %-40s  %s\n", "TARGET", "STATUS", "SOURCE", "DEPLOY DIR")
	fmt.Printf("  %-12s  %-10s  %-40s  %s\n", "------------", "----------", "----------------------------------------", "----------")
	for _, row := range rows {
		sourceCol := row.Source
		if row.Status == "inactive" && row.Needs != "" {
			sourceCol = "needs " + row.Needs
		}
		fmt.Printf("  %-12s  %-10s  %-40s  %s\n", row.Target, row.Status, sourceCol, row.DeployDir)
	}

	hasActive := false
	for _, r := range rows {
		if r.Status == "active" {
			hasActive = true
			break
		}
	}
	if !hasActive {
		fmt.Println()
		fmt.Println("[i] Create a harness config (e.g. CLAUDE.md, .cursor/, .github/copilot-instructions.md)")
		fmt.Println("    or declare `targets:` in apm.yml.")
	}
	return nil
}

// findFile checks if path exists relative to root.
func findFile(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, rel))
	return err == nil
}

var _ = findFile // suppress unused warning
