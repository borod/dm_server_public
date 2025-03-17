package dm_configuration

type Configuration struct {
	Name   string       `json:"Name"`
	DB     DatabaseConf `json:"DB"`
	Crypto CryptoConf   `json:"Crypto"`
	Mail   MailConf     `json:"Mail"`
}

type DBConf struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DBName   string `json:"DBName"`
}

type DatabaseConf struct {
	MySQL DBConf `json:"MySQL"`
	Mongo DBConf `json:"Mongo"`
}

type MailConf struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	From     string `Json:"From"`
}

type CryptoConf struct {
	HashFunction string `json:"HashFunction"`
}
