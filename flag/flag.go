package flag

import (
	"strings"

	goflag "flag"

	"github.com/spf13/pflag"
)

type NamedFlagSet struct {
	Order []string

	FlagSets map[string]*pflag.FlagSet
}

func (nfs *NamedFlagSet) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}

	if _, ok := nfs.FlagSets[name]; !ok {
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}
	return nfs.FlagSets[name]
}

func InitFlags() {
	pflag.CommandLine.SetNormalizeFunc(WordNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func WordNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}

	return pflag.NormalizedName(name)
}
