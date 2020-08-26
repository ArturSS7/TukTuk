package config

//Initialization structures

type InitStruct struct {
	TelegramBot      TelegramSetting
	AdminCredentials AdminPanelCredentials
	DomainConfig     Domain
	EmailAlert       EmailAlertSetting
	HttpsCertPath    HttpsConfig
}

type TelegramSetting struct {
	Token       string `json:"token"`
	ChatID      int64  `json:"chat_id"`
	LengthAlert string `json:"length_alert"`
	Enabled     bool   `json:"enabled"`
}

type Domain struct {
	Name             string `json:"name"`
	IPV4             string `json:"ipv4"`
	NonExistingIPV4  string `json:"non_existing_ipv4"`
	IPV6             string `json:"ipv6"`
	NonExistingIPV6  string `json:"non_existing_ipv6"`
	AcmeTxtChallenge string `json:"acme_txt_challenge"`
}

type AdminPanelCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EmailAlertSetting struct {
	To      string `json:"to"`
	Enabled bool   `json:"enabled"`
}

type HttpsConfig struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
