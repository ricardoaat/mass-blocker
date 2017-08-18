package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

/*LoadConfig Loads the configurations parameters stored on
* the configuration file (config.toml)
 */
func LoadConfig(configFile string) error {

	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal("Load config failed", err)
		return err
	}
	return nil
}
