package startinitialization

//Initialization structures

type InitStruct struct {
	Telegrambot TelegramSetting
	Admincred   AdminPanelCredentials
	Domain      string
	EmailAlert  EmailAlertSetting
}

type TelegramSetting struct {
	Token       string
	ChatID      int64
	LenghtAlert string
	Enabled     bool
}
type AdminPanelCredentials struct {
	Username string
	Password string
}

type EmailAlertSetting struct {
	To      string
	Enabled bool
}
