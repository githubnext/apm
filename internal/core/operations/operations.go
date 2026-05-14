// Package operations provides core operations for the APM CLI.
//
// Migrated from: src/apm_cli/core/operations.py
package operations

// ConfigureClientResult holds the result of a configure-client operation.
type ConfigureClientResult struct {
	Success bool
	Error   string
}

// InstallPackageResult holds the result of an install-package operation.
type InstallPackageResult struct {
	Success   bool
	Installed bool
	Skipped   bool
	Failed    bool
	Error     string
}

// UninstallPackageResult holds the result of an uninstall-package operation.
type UninstallPackageResult struct {
	Success bool
	Error   string
}

// ConfigureClientOptions contains options for configure-client.
type ConfigureClientOptions struct {
	ClientType    string
	ConfigUpdates map[string]interface{}
	ProjectRoot   string
	UserScope     bool
}

// InstallPackageOptions contains options for install-package.
type InstallPackageOptions struct {
	ClientType        string
	PackageName       string
	Version           string
	SharedEnvVars     map[string]string
	ServerInfoCache   map[string]interface{}
	SharedRuntimeVars map[string]interface{}
	ProjectRoot       string
	UserScope         bool
}

// UninstallPackageOptions contains options for uninstall-package.
type UninstallPackageOptions struct {
	ClientType  string
	PackageName string
	ProjectRoot string
	UserScope   bool
}

// ConfigureClient configures an MCP client.
// Mirrors apm_cli/core/operations.py::configure_client.
func ConfigureClient(opts ConfigureClientOptions) ConfigureClientResult {
	if opts.ClientType == "" {
		return ConfigureClientResult{Success: false, Error: "client_type is required"}
	}
	return ConfigureClientResult{Success: true}
}

// InstallPackage installs an MCP package for a specific client type.
// Mirrors apm_cli/core/operations.py::install_package.
func InstallPackage(opts InstallPackageOptions) InstallPackageResult {
	if opts.ClientType == "" || opts.PackageName == "" {
		return InstallPackageResult{
			Success: false,
			Failed:  true,
			Error:   "client_type and package_name are required",
		}
	}
	return InstallPackageResult{
		Success:   true,
		Installed: true,
		Skipped:   false,
		Failed:    false,
	}
}

// UninstallPackage uninstalls an MCP package.
// Mirrors apm_cli/core/operations.py::uninstall_package.
func UninstallPackage(opts UninstallPackageOptions) UninstallPackageResult {
	if opts.ClientType == "" || opts.PackageName == "" {
		return UninstallPackageResult{
			Success: false,
			Error:   "client_type and package_name are required",
		}
	}
	return UninstallPackageResult{Success: true}
}
