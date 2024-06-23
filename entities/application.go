package entities

import "fmt"

const (
	PRINTOUT_TYPE_LOG = iota
	PRINTOUT_TYPE_ERR
)

var (
	PrintOutChan = make(chan PrintOut)

	//* WS Requirement
	CentralAuthHost  = "127.0.0.1:8081" //! Buat diharcode, jangan konfigurasi
	SessionMode      = true             //! buat hardcode, jangan konfigurasi
	SecureWSProtocol = false            //! Isi false jika tidak menggunakan ssl, isi true jika menggunakan ssl
)

type PrintOut struct {
	Type    int
	Message []interface{}
}

func PrintErrorf(format string, a ...any) {
	PrintError(fmt.Sprintf(format, a...))
}

func PrintLogf(format string, a ...any) {
	PrintLog(fmt.Sprintf(format, a...))
}

func PrintError(message ...interface{}) {
	po := PrintOut{
		Type:    PRINTOUT_TYPE_ERR,
		Message: message,
	}

	PrintOutChan <- po
}

func PrintLog(message ...interface{}) {
	po := PrintOut{
		Type:    PRINTOUT_TYPE_LOG,
		Message: message,
	}
	PrintOutChan <- po
}
