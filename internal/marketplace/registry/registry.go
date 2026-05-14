// Package registry manages registered marketplaces stored in ~/.apm/marketplaces.json.
package registry

import (
"encoding/json"
"fmt"
"os"
"path/filepath"
"sort"
"strings"
"sync"
)

const marketplacesFilename = "marketplaces.json"

// MarketplaceSource represents a registered marketplace.
type MarketplaceSource struct {
Name string `json:"name"`
URL  string `json:"url"`
// Additional fields are preserved via Extra.
Extra map[string]interface{} `json:"-"`
}

// FromDict creates a MarketplaceSource from a JSON-decoded map.
func FromDict(m map[string]interface{}) (MarketplaceSource, error) {
name, ok := m["name"].(string)
if !ok || name == "" {
return MarketplaceSource{}, fmt.Errorf("missing or invalid 'name' field")
}
url, _ := m["url"].(string)
extra := make(map[string]interface{})
for k, v := range m {
if k != "name" && k != "url" {
extra[k] = v
}
}
return MarketplaceSource{Name: name, URL: url, Extra: extra}, nil
}

// ToDict converts a MarketplaceSource to a JSON-serializable map.
func (s MarketplaceSource) ToDict() map[string]interface{} {
m := make(map[string]interface{}, len(s.Extra)+2)
for k, v := range s.Extra {
m[k] = v
}
m["name"] = s.Name
m["url"] = s.URL
return m
}

// Registry manages the marketplace list file.
type Registry struct {
configDir func() string
mu        sync.Mutex
cache     []MarketplaceSource
cacheValid bool
}

// New creates a Registry that stores files in the directory returned by configDir.
func New(configDir func() string) *Registry {
return &Registry{configDir: configDir}
}

func (r *Registry) path() string {
return filepath.Join(r.configDir(), marketplacesFilename)
}

func (r *Registry) ensureFile() (string, error) {
dir := r.configDir()
if err := os.MkdirAll(dir, 0o755); err != nil {
return "", err
}
p := r.path()
if _, err := os.Stat(p); os.IsNotExist(err) {
data, _ := json.MarshalIndent(map[string]interface{}{"marketplaces": []interface{}{}}, "", "  ")
if err := os.WriteFile(p, data, 0o644); err != nil {
return "", err
}
}
return p, nil
}

func (r *Registry) invalidate() {
r.mu.Lock()
r.cacheValid = false
r.mu.Unlock()
}

func (r *Registry) load() ([]MarketplaceSource, error) {
r.mu.Lock()
defer r.mu.Unlock()
if r.cacheValid {
out := make([]MarketplaceSource, len(r.cache))
copy(out, r.cache)
return out, nil
}
p, err := r.ensureFile()
if err != nil {
return nil, err
}
raw, err := os.ReadFile(p)
var data map[string]interface{}
if err == nil {
_ = json.Unmarshal(raw, &data)
}
if data == nil {
data = map[string]interface{}{"marketplaces": []interface{}{}}
}
entries, _ := data["marketplaces"].([]interface{})
var sources []MarketplaceSource
for _, e := range entries {
m, ok := e.(map[string]interface{})
if !ok {
continue
}
src, err := FromDict(m)
if err == nil {
sources = append(sources, src)
}
}
r.cache = sources
r.cacheValid = true
out := make([]MarketplaceSource, len(sources))
copy(out, sources)
return out, nil
}

func (r *Registry) save(sources []MarketplaceSource) error {
p, err := r.ensureFile()
if err != nil {
return err
}
dicts := make([]interface{}, len(sources))
for i, s := range sources {
dicts[i] = s.ToDict()
}
data := map[string]interface{}{"marketplaces": dicts}
raw, err := json.MarshalIndent(data, "", "  ")
if err != nil {
return err
}
tmp := p + ".tmp"
if err := os.WriteFile(tmp, raw, 0o644); err != nil {
return err
}
if err := os.Rename(tmp, p); err != nil {
return err
}
r.mu.Lock()
r.cache = make([]MarketplaceSource, len(sources))
copy(r.cache, sources)
r.cacheValid = true
r.mu.Unlock()
return nil
}

// GetAll returns all registered marketplaces.
func (r *Registry) GetAll() ([]MarketplaceSource, error) {
return r.load()
}

// GetByName returns a marketplace by display name (case-insensitive).
// Returns an error if not found.
func (r *Registry) GetByName(name string) (MarketplaceSource, error) {
lower := strings.ToLower(name)
sources, err := r.load()
if err != nil {
return MarketplaceSource{}, err
}
for _, s := range sources {
if strings.ToLower(s.Name) == lower {
return s, nil
}
}
return MarketplaceSource{}, fmt.Errorf("marketplace not found: %s", name)
}

// Add registers a marketplace, replacing any existing entry with the same name.
func (r *Registry) Add(source MarketplaceSource) error {
sources, err := r.load()
if err != nil {
return err
}
lower := strings.ToLower(source.Name)
var filtered []MarketplaceSource
for _, s := range sources {
if strings.ToLower(s.Name) != lower {
filtered = append(filtered, s)
}
}
filtered = append(filtered, source)
return r.save(filtered)
}

// Remove removes a marketplace by name.
// Returns an error if not found.
func (r *Registry) Remove(name string) error {
sources, err := r.load()
if err != nil {
return err
}
lower := strings.ToLower(name)
var filtered []MarketplaceSource
for _, s := range sources {
if strings.ToLower(s.Name) != lower {
filtered = append(filtered, s)
}
}
if len(filtered) == len(sources) {
return fmt.Errorf("marketplace not found: %s", name)
}
return r.save(filtered)
}

// Names returns a sorted list of registered marketplace names.
func (r *Registry) Names() ([]string, error) {
sources, err := r.load()
if err != nil {
return nil, err
}
names := make([]string, len(sources))
for i, s := range sources {
names[i] = s.Name
}
sort.Strings(names)
return names, nil
}

// Count returns the number of registered marketplaces.
func (r *Registry) Count() (int, error) {
sources, err := r.load()
if err != nil {
return 0, err
}
return len(sources), nil
}
