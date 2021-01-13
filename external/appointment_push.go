package external

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

var appointmentPushUrl = ""

type PushAppointmentRequest struct {
	*echoapp.Appointment
}

type PushAppointmentResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 推送预约信息到旅游局
func DoPushAppointmentRequest(request *PushAppointmentRequest) (*PushAppointmentResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest marshal")
	}

	res, err := http.Post(appointmentPushUrl, echo.MIMEApplicationJSON, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest post")
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest status code is not 200")
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest readAll")
	}
	var response PushAppointmentResponse
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest unmarshal")
	}

	if response.Code != 0 {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest errMsg "+response.Msg)
	}
	return &response, nil
}
