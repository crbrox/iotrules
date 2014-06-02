package config

import (
	"fmt"
	"iotrules/mylog"
)

var port string
var (
	smsEndPoint string
	apiSecret   string
	apiKey      string
)
var (
	smtpServer string
)

var (
	updateEndPoint string
)

var loaded bool = false

func LoadConfig(filename string) (err error) {
	if mylog.Debugging {
		mylog.Debugf("enter config.LoadConfig %q", filename)
		defer func() { mylog.Debugf("exit config.LoadConfig %+v", err) }()
	}

	err = loadFile(filename)
	if err != nil {
		return err
	}

	var found bool

	port, found = str("port")
	if !found {
		return fmt.Errorf("config port is mandatory")
	}

	smsEndPoint, found = str("SMS.endpoint")
	if !found {
		return fmt.Errorf("config SMS.endpoint is mandatory")
	}

	apiKey, found = str("API_KEY")
	if !found {
		return fmt.Errorf("config SMS.endpoint is mandatory")
	}

	apiSecret, found = str("API_SECRET")
	if !found {
		return fmt.Errorf("config SMS.endpoint is mandatory")
	}

	smtpServer, found = str("email.SMTP.server")
	if !found {
		return fmt.Errorf("config SMS.endpoint is mandatory")
	}

	updateEndPoint, found = str("update.endpoint")
	if !found {
		return fmt.Errorf("config update.endpoint is mandatory")
	}

	return nil
}

func Port() string           { return port }
func SMSEndpoint() string    { return smsEndPoint }
func APIKey() string         { return apiKey }
func APISecret() string      { return apiSecret }
func SMTPServer() string     { return smtpServer }
func UpdateEndpoint() string { return updateEndPoint }
