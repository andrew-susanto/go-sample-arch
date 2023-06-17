package config

// Config is a struct that contains all configuration
type Config struct {
	FeatureFlag     map[string]interface{}      `yaml:"feature_flag"`
	AWS             ConfigAWS                   `yaml:"aws"`
	HttpClient      map[string]ConfigHttpClient `yaml:"http_client"`
	RpcClientConfig RpcClientMap                `yaml:"rpc_client"`
	SQSConfig       SQSConfigMap                `yaml:"sqs"`
}

// ConfigAWS is a struct that contains all configuration for AWS
type ConfigAWS struct {
	Region string `yaml:"region"`
}

// ConfigHttpClient is a struct that contains all configuration for HTTP Client
type ConfigHttpClient struct {
	Timeout                  int `yaml:"timeout"`
	MaxConcurrentRequest     int `yaml:"max_concurrent"`
	ErrorPercentageThreshold int `yaml:"error_percentage"`
}

// RpcClientMap is a struct that contains all configuration for RPC Client
type RpcClientMap struct {
	CxpIgwTrip RpcClientConfig `yaml:"cxpigw_trip"`
}

// RpcClientConfig is a struct that contains all configuration for RPC Client
type RpcClientConfig struct {
	ServiceUrl               string `yaml:"service_url"`
	ServiceName              string `yaml:"service_name"`
	Timeout                  int    `yaml:"timeout"`
	MaxConcurrentRequest     int    `yaml:"max_concurrent"`
	ErrorPercentageThreshold int    `yaml:"error_percentage"`
}

// SQSConfigMap is a struct that contains all configuration for SQS
type SQSConfigMap struct {
	IssuanceJobFIFO SQSClientConfig `yaml:"issuance_job_fifo"`
	MyInbox         SQSClientConfig `yaml:"myinbox"`
}

// SQSClientConfig is a struct that contains all configuration for SQS
type SQSClientConfig struct {
	QueueName              string `yaml:"queue_name"`
	MaxNumberMessage       int    `yaml:"max_number_message"`
	PollPeriodInMilisecond int    `yaml:"poll_period_in_milisecond"`
	Enabled                bool   `yaml:"enabled"`
}
