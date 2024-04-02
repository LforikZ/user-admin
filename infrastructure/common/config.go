package common

type Mconfig struct {
	ConfigMyResponse map[int]string
}

var (
	MyConfig Mconfig
)

func init() {
	MyConfig.ConfigMyResponse = getResponseConfig()
}
