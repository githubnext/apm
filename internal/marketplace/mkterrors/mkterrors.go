// Package mkterrors defines the marketplace error hierarchy.
package mkterrors

import "fmt"

// MarketplaceError is the base type for marketplace errors.
type MarketplaceError struct {
msg string
}

func (e *MarketplaceError) Error() string { return e.msg }

// MarketplaceNotFoundError is raised when a marketplace cannot be found.
type MarketplaceNotFoundError struct {
Name string
Host string
MarketplaceError
}

// NewMarketplaceNotFoundError creates a MarketplaceNotFoundError.
func NewMarketplaceNotFoundError(name, host string) *MarketplaceNotFoundError {
if host == "" {
host = "github.com"
}
return &MarketplaceNotFoundError{
Name: name,
Host: host,
MarketplaceError: MarketplaceError{
msg: fmt.Sprintf("Marketplace '%s' is not registered. Run 'apm marketplace add https://%s/OWNER/REPO' to register it.", name, host),
},
}
}

// PluginNotFoundError is raised when a plugin is not found.
type PluginNotFoundError struct {
PluginName      string
MarketplaceName string
MarketplaceError
}

// NewPluginNotFoundError creates a PluginNotFoundError.
func NewPluginNotFoundError(pluginName, marketplaceName string) *PluginNotFoundError {
return &PluginNotFoundError{
PluginName:      pluginName,
MarketplaceName: marketplaceName,
MarketplaceError: MarketplaceError{
msg: fmt.Sprintf("Plugin '%s' not found in marketplace '%s'.", pluginName, marketplaceName),
},
}
}

// MarketplaceYmlError is raised when marketplace.yml validation fails.
type MarketplaceYmlError struct {
Message string
MarketplaceError
}

// NewMarketplaceYmlError creates a MarketplaceYmlError.
func NewMarketplaceYmlError(message string) *MarketplaceYmlError {
return &MarketplaceYmlError{Message: message, MarketplaceError: MarketplaceError{msg: message}}
}

// MarketplaceFetchError is raised when fetching marketplace data fails.
type MarketplaceFetchError struct {
MarketplaceError
}

// NewMarketplaceFetchError creates a MarketplaceFetchError.
func NewMarketplaceFetchError(msg string) *MarketplaceFetchError {
return &MarketplaceFetchError{MarketplaceError: MarketplaceError{msg: msg}}
}
