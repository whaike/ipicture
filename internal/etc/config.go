package etc

import (
	"github.com/pyroscope-io/pyroscope/pkg/util/file"
	"ipicture/g"
	"os"
	"sigs.k8s.io/yaml"
)

type Config struct {
	ZapLog          g.ZapLogConf
	Path            string `yaml:"Path"`
	PyroscopeEnable bool   `yaml:"PyroscopeEnable"`
	PyroscopeAddr   string `yaml:"PyroscopeAddr"`
	DelDuplicate    bool   `yaml:"DelDuplicate"`
}

func LoadConfig(filepath string) *Config {
	if !file.Exists(filepath) {
		return &Config{
			ZapLog:          g.ZapLogConf{},
			Path:            ".",
			PyroscopeEnable: false,
			PyroscopeAddr:   "",
		}
	}
	f, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	c := &Config{}
	if err = yaml.Unmarshal(f, c); err != nil {
		panic(err)
	}
	return c
}
