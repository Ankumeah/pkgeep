package modes

var configFileFlag = flag{
	name:        "--config-file",
	shortName:   "-c",
	description: "Path to the config file",
	usage:       "(--config-file,-c)=path/to/config/file",
}
