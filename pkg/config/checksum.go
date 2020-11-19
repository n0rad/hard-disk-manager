package config

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/integrity"
	"regexp"
)

type ChecksumConfig struct {
	Pattern            string
	PatternIsExclusive bool
	Hash               hashs.Hash
	Strategy           string

	regex *regexp.Regexp
}

func (h *ChecksumConfig) Init() error {
	if h.Pattern == "" {
		h.Pattern = `(?i)\.*$`
	}
	var err error
	h.regex, err = regexp.Compile(h.Pattern)
	if err != nil {
		return errs.WithEF(err, data.WithField("regex", h.regex), "Failed to compile files regex")
	}

	if h.Hash == "" {
		h.Hash = hashs.Sha256
	}

	if h.Strategy == "" {
		h.Strategy = integrity.SumfileStrategy
	}

	return nil
}

func (h ChecksumConfig) GetRegex() *regexp.Regexp {
	return h.regex
}
