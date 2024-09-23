package usecase

import (
	"central-dummy-master/entities"
	"central-dummy-master/repository/sessionrepo"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kpango/glg"
	"github.com/randyardiansyah25/libpkg/util/env"
	"github.com/randyardiansyah25/wsbase-handler"
)

// * CA Requirement
func OpenCentralSession() (er error) {

	//* Generate Central Code
	//* Format : sha512(<database name>|<database host>|<port>)
	toHash := fmt.Sprintf("%s|%s|%s", env.GetString("postgres.dbname"), env.GetString("postgres.host"), env.GetString("postgres.port"))
	cryptor := sha512.New()
	cryptor.Write([]byte(toHash))
	centralCode := hex.EncodeToString(cryptor.Sum(nil))

	repo := sessionrepo.NewSessionRepo()
	sessionId, er := repo.OpenSession(centralCode)
	if er != nil {
		return
	}

	//* Check: jika menggunakan session mode dari pengaturan, maka jalankan websocket client
	if entities.SessionMode {
		//* format: /session/start/<session id>/<client id>/<central_code>/<api key>
		wsPath := fmt.Sprintf("/session/start/%s/%s/%s/%s",
			sessionId,
			env.GetString("app.client_id"),
			centralCode,
			env.GetString("app.api_key"),
		)

		client := wsbase.NewWSClient(entities.CentralAuthHost, wsPath, entities.SecureWSProtocol)
		client.SetLogHandler(func(logType int, val string) {
			switch logType {
			case wsbase.LOG:
				_ = glg.Log(val)
			case wsbase.ERR:
				_ = glg.Error(val)
			default:
				_ = glg.Info(val)
			}
		})

		period := env.GetInt("app.reconnect_period")
		client.SetReconnectPeriod(time.Duration(period) * time.Second)

		//* Digunakan untuk penanganan saat menerima pesan dari websocket server
		//* Lebih lanjut baca di https://github.com/randyardiansyah25/wsbase-handler/blob/master/README.md#client
		client.SetMessageHandler(func(m wsbase.Message) {
			j, _ := json.MarshalIndent(m, "", "    ")
			entities.PrintLog("receive message : ")
			fmt.Println(string(j))
		})

		go client.Start()
		er = <-wsbase.WSClientErrorSignal
	}

	return
}
