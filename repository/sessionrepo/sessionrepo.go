package sessionrepo

import (
	"central-dummy-master/entities"
	"encoding/json"
	"fmt"
	"net/http"

	net "github.com/randyardiansyah25/libpkg/net/http"
	"github.com/randyardiansyah25/libpkg/util/env"
)

// * CA Requirement
type SessionRepo interface {
	OpenSession(centralCode string) (sessionId string, er error)
}

func NewSessionRepo() SessionRepo {
	return &sessionImpl{}
}

type sessionImpl struct {
}

func (s sessionImpl) OpenSession(centralCode string) (sessionId string, er error) {
	caUrl := fmt.Sprintf("http://%s/session/open", entities.CentralAuthHost)
	client := net.NewSimpleClient("POST", caUrl, 30)
	client.SetAuthorization(env.GetString("app.api_key"))
	client.SetContentTypeFormUrlEncoded()
	client.SetHeader("Client-ID", env.GetString("app.client_id"))

	client.AddParam("central_code", centralCode)

	resp := struct {
		ResponseCode    string `json:"response_code"`
		ResponseMessage string `json:"response_message"`
		ResponseData    struct {
			RequestData any    `json:"request_data"`
			SessionId   string `json:"session_id"`
		} `json:"response_data"`
	}{}

	entities.PrintLog("Authorizing central...")
	res := client.DoRequest()
	if res.StatusCode() == http.StatusOK {
		if er = json.Unmarshal([]byte(res.Message()), &resp); er != nil {
			er = fmt.Errorf("error unmarshal central authorization response: %s", er.Error())
		} else {
			if resp.ResponseCode == "00" {
				sessionId = resp.ResponseData.SessionId
			} else {
				er = fmt.Errorf("[%s] authorize failed: %s", resp.ResponseCode, resp.ResponseMessage)
			}
		}

	} else {
		er = fmt.Errorf("authorize failed: [%d] %s", res.StatusCode(), res.Message())
	}
	return
}
