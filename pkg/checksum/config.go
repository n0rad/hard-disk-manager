package checksum

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"
	"io/ioutil"
	"regexp"
)

type Config struct {
	Pattern            string
	PatternIsInclusive bool
	Hash               hashs.Hash
	Strategy           string

	regex *regexp.Regexp
}

func (h *Config) Init() error {
	if h.Pattern == "" {
		h.Pattern = `(?i)\.*$`
	}

	var err error
	h.regex, err = regexp.Compile(h.Pattern)
	if err != nil {
		return errs.WithEF(err, data.WithField("Regex", h.regex), "Failed to compile files Regex")
	}

	return nil
}

func (h *Config) Load(configPath string) error {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errs.WithEF(err, data.WithField("path", configPath), "Failed to read fim config file")
	}

	if err := yaml.Unmarshal(bytes, h); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)).WithField("path", configPath), "Failed to parse fim file")
	}

	if err := h.Init(); err != nil {
		return errs.WithEF(err, data.WithField("content", string(bytes)), "Failed to init hdm file")
	}
	return nil
}
