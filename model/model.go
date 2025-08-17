package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

type User struct {
	Name     string
	Age      int
	Email    string
	Password string
	Contact  string
	Address  Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

// Logger interface to support custom logging libraries
type Logger interface {
	Fatal(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
	Info(string, ...any)
	Debug(string, ...any)
	Trace(string, ...any)
}

// Driver is our database driver
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
	log     Logger
}

// Options for creating a new driver
type Options struct {
	Logger Logger
}

// Stat checks file/directory existence with support for `.json` extension
func Stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

// New creates a new database driver
func New(dir string, opts *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	// ❌ Bug: you shadowed the opts parameter before
	// ✅ Fix: copy properly
	var options Options
	if opts != nil {
		options = *opts
	}

	if options.Logger == nil {
		options.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     options.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		options.Logger.Debug("Using '%s' (Database already exists)\n", dir)
		return &driver, nil
	}

	options.Logger.Debug("Creating database at '%s'\n", dir)
	return &driver, os.Mkdir(dir, 0755)
}

// Write stores a resource inside a collection
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing Collection - no place to save records")
	}
	if resource == "" {
		return fmt.Errorf("Missing Resource - no place to save records (no file name)")
	}

	mutex := d.GetOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)

	// ❌ Bug: you used os.Mkdir instead of MkdirAll
	// ✅ Fix: MkdirAll ensures dir exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	b, err := json.MarshalIndent(v, "", "\t ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

// Read loads a resource from a collection and returns JSON string
func (d *Driver) Read(collection, resource string) (string, error) {
	if collection == "" {
		return "", fmt.Errorf("Missing Collection - no place to save records")
	}
	if resource == "" {
		return "", fmt.Errorf("Missing Resource - no place to save records (no file name)")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := Stat(record); err != nil {
		return "", err
	}

	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ReadAll loads all resources in a collection and returns them as JSON strings
func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing Collection - no place to save records")
	}
	dir := filepath.Join(d.dir, collection)

	if _, err := Stat(dir); err != nil {
		return nil, err
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

// Delete removes a resource from a collection
func (d *Driver) Delete(collection, resource string) error {
	if collection == "" {
		return fmt.Errorf("Missing Collection - no place to save records")
	}
	if resource == "" {
		return fmt.Errorf("Missing Resource - no place to save records (no file name)")
	}

	record := filepath.Join(d.dir, collection, resource+".json")

	mutex := d.GetOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	// ❌ Bug: you joined path twice
	// ✅ Fix: just use record directly
	switch fi, err := Stat(record); {
	case err != nil:
		return err
	case fi.Mode().IsDir():
		return os.RemoveAll(record)
	case fi.Mode().IsRegular():
		return os.Remove(record)
	default:
		return os.Remove(record)
	}
}

// GetOrCreateMutex provides per-collection locking
func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}
