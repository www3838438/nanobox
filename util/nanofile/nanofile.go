package nanofile

import (
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/nanobox-io/nanobox/util"
)

var vip *viper.Viper

func Viper() *viper.Viper {
	if vip != nil {
		return vip
	}

	vip = viper.New()
	vip.SetDefault("cpu-cap",  50)
	vip.SetDefault("cpus",  2)
	vip.SetDefault("host-dns", "off")
	vip.SetDefault("mount-nfs",  true)
	vip.SetDefault("name",  util.LocalDirName())
	vip.SetDefault("provider", "virtualbox") // this may change in the future (adding additional hosts such as vmware
	vip.SetDefault("ram",  1024)
	vip.SetDefault("use-proxy",  false)

	vip.SetConfigFile(filepath.Join(util.GlobalDir(), "nanofile.yml"))
	vip.MergeInConfig() // using merge because it starts from existing config
	vip.SetConfigFile(filepath.Join(util.LocalDir(), "nanofile.yml"))
	vip.MergeInConfig()
	return vip
}

