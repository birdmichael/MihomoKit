package mihomo

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
)

type RealTimeLogger interface {
	Log(level string, payload string)
}

var (
	logger       RealTimeLogger
	runtimeMutex sync.Mutex
	started      bool
	configCache  []byte
	configPath   string
	loggerOnce   sync.Once
)

// Setup initialises the Mihomo runtime using the provided working directory and configuration.
//
// The configSource parameter can be either an absolute/relative file path, a base64 encoded YAML string
// or a raw YAML string. The resolved configuration is persisted to the default Mihomo config path
// so that subsequent reloads can operate on the same source.
func Setup(homeDir string, configSource string) {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	loggerOnce.Do(func() {
		go SubscribeLogger()
	})

	if err := initialiseHome(homeDir); err != nil {
		log.Errorln("initialise home dir failed: %s", err)
		return
	}

	cfgBytes, resolvedPath, err := resolveConfig(configSource)
	if err != nil {
		log.Errorln("load config failed: %s", err)
		return
	}

	if err := writeConfigToDisk(cfgBytes, resolvedPath); err != nil {
		log.Warnln("persist config warning: %s", err)
	}

	if err := config.Init(C.Path.HomeDir()); err != nil {
		log.Warnln("initialise config directory warning: %s", err)
	}

	ensureGeoDatasets()

	if err := applyConfigLocked(cfgBytes); err != nil {
		log.Errorln("apply config failed: %s", err)
		return
	}

	configCache = append(configCache[:0], cfgBytes...)
	configPath = resolvedPath
	started = true
}

// Reload re-applies the last successfully loaded configuration.
func Reload() error {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if !started {
		return errors.New("mihomo has not been started")
	}

	return applyConfigLocked(configCache)
}

// ReloadWithConfig re-applies Mihomo with the provided raw configuration bytes.
func ReloadWithConfig(data []byte) error {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if len(data) == 0 {
		return errors.New("configuration payload is empty")
	}

	if err := applyConfigLocked(data); err != nil {
		return err
	}

	configCache = append(configCache[:0], data...)
	if err := writeConfigToDisk(data, configPath); err != nil {
		log.Warnln("persist config warning: %s", err)
	}
	return nil
}

// ReloadWithBase64 decodes the base64 configuration payload prior to reloading.
func ReloadWithBase64(configBase64 string) error {
	buf, err := base64.StdEncoding.DecodeString(configBase64)
	if err != nil {
		return fmt.Errorf("decode config failed: %w", err)
	}
	return ReloadWithConfig(buf)
}

// Stop gracefully shuts down running Mihomo services.
func Stop() {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if !started {
		return
	}

	executor.Shutdown()
	started = false
}

func SetupLogger(l RealTimeLogger) {
	logger = l
}

func SubscribeLogger() {
	sub := log.Subscribe()
	defer log.UnSubscribe(sub)

	for ev := range sub {
		if logger != nil {
			entry := ev
			logger.Log(entry.Type(), entry.Payload)
		}
	}
}

func initialiseHome(homeDir string) error {
	if homeDir == "" {
		return errors.New("homeDir is empty")
	}

	absHome := homeDir
	if !filepath.IsAbs(homeDir) {
		var err error
		absHome, err = filepath.Abs(homeDir)
		if err != nil {
			return fmt.Errorf("resolve home dir failed: %w", err)
		}
	}

	if err := os.MkdirAll(absHome, 0o755); err != nil {
		return fmt.Errorf("create home dir failed: %w", err)
	}

	C.SetHomeDir(absHome)
	return nil
}

func resolveConfig(input string) ([]byte, string, error) {
	trimmed := strings.TrimSpace(input)

	if trimmed != "" {
		resolvedPath := trimmed
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Join(C.Path.HomeDir(), resolvedPath)
		}
		if stat, err := os.Stat(resolvedPath); err == nil && !stat.IsDir() {
			data, err := os.ReadFile(resolvedPath)
			if err != nil {
				return nil, "", fmt.Errorf("read config file failed: %w", err)
			}
			return data, resolvedPath, nil
		}
	}

	if decoded, err := base64.StdEncoding.DecodeString(trimmed); err == nil {
		return decoded, filepath.Join(C.Path.HomeDir(), C.Path.Config()), nil
	}

	if trimmed == "" {
		return nil, "", errors.New("configuration is empty")
	}
	return []byte(trimmed), filepath.Join(C.Path.HomeDir(), C.Path.Config()), nil
}

func ensureGeoDatasets() {
	paths := []string{
		C.Path.MMDB(),
		C.Path.ASN(),
		C.Path.GeoSite(),
		C.Path.GeoIP(),
	}

	for _, path := range paths {
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
	}
}

func writeConfigToDisk(data []byte, target string) error {
	if len(data) == 0 {
		return errors.New("configuration payload is empty")
	}

	if target == "" {
		target = filepath.Join(C.Path.HomeDir(), C.Path.Config())
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("prepare config directory failed: %w", err)
	}

	if err := os.WriteFile(target, data, 0o600); err != nil {
		return fmt.Errorf("write config failed: %w", err)
	}

	C.SetConfig(target)
	return nil
}

func applyConfigLocked(data []byte) error {
	if len(data) == 0 {
		return errors.New("configuration payload is empty")
	}

	if err := hub.Parse(data); err != nil {
		return err
	}

	return nil
}
