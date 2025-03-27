package configs

type ConfigManagerConfig struct {
	ServiceName         string `name:"service_name" type:"string" description:"service name" required:"true"`
	ServiceDescription  string `name:"service_description" type:"bigText" description:"This config is to update the description of this service" required:"false"`
	ConfigRefreshInSecs int    `name:"config_refresh_in_secs" type:"number" description:"config refresh in seconds" required:"false"`
}

func (c *ConfigManagerConfig) WithDefault() {
	if c.ServiceName == "" {
		c.ServiceName = "Config Manager Service"
	}
}

func (c *ConfigManagerConfig) Key() string {
	return "config-manager-config"
}
