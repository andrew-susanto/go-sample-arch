package secretmanager

// Secrets is a struct that contains all secrets
type Secrets struct {
	Postgresql SecretsPostgreSQL `json:"postgresql"`
	Redis      SecretsRedis      `json:"redis"`
	DocDB      SecretsDocDb      `json:"docdb"`
	OpenSearch SecretsOpenSearch `json:"opensearch"`
}

// SecretsPostgreSQL is a struct that contains all postgresql secrets
type SecretsPostgreSQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

// SecretsRedis is a struct that contains all redis secrets
type SecretsRedis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// SecretsDocDb is a struct that contains all docdb secrets
type SecretsDocDb struct {
	DBName          string `json:"db_name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ClusterEndpoint string `json:"clusterendpoint"`
}

// SecretsOpenSearch is a struct that contains all opensearch secrets
type SecretsOpenSearch struct {
	Domain    string   `json:"domain"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Addresses []string `json:"addresses"`
}
