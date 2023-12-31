package am

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

type (
	Core interface {
		Name() string
		Log() Logger
		Cfg() *Config
		Lifecycle
	}

	Lifecycle interface {
		Setup(ctx context.Context) error
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}
)

type (
	SimpleCore struct {
		name     string
		log      Logger
		cfg      *Config
		didSetup bool
		didStart bool
	}
)

func NewCore(name string, opts ...Option) *SimpleCore {
	name = GenName(name, "core")

	bw := &SimpleCore{
		name: name,
	}

	for _, opt := range opts {
		opt(bw)
	}

	return bw
}

func (sc *SimpleCore) Name() string {
	return sc.name
}

func (sc *SimpleCore) SetName(name string) {
	sc.name = name
}

func (sc *SimpleCore) Log() Logger {
	return sc.log
}

func (sc *SimpleCore) SetLog(log Logger) {
	sc.log = log
}

func (sc *SimpleCore) Cfg() *Config {
	return sc.cfg
}

func (sc *SimpleCore) SetCfg(cfg *Config) {
	sc.cfg = cfg
}

func (sc *SimpleCore) Setup(ctx context.Context) error {
	sc.Log().Infof("%s setup", sc.Name())
	return nil
}

func (sc *SimpleCore) Start(ctx context.Context) error {
	sc.Log().Infof("%s start", sc.Name())
	return nil
}

func (sc *SimpleCore) Stop(ctx context.Context) error {
	sc.Log().Infof("%s stop", sc.Name())
	return nil
}

func GenName(name, defName string) string {
	if strings.Trim(name, " ") == "" {
		return fmt.Sprintf("%s-%s", defName, nameSufix())
	}
	return name
}

func nameSufix() string {
	digest := hash(time.Now().String())
	return digest[len(digest)-8:]
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

type (
	Option func(w *SimpleCore)
)

func WithConfig(cfg *Config) Option {
	return func(sc *SimpleCore) {
		sc.SetCfg(cfg)
	}
}

func WithLogger(log Logger) Option {
	return func(sc *SimpleCore) {
		sc.SetLog(log)
	}
}
