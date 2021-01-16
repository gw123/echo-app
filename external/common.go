package external

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

var app_key = "mZ0z1Df0BirPWnrvGVBg5FKqR_B-uDVOnIlTQbofexQ"

func DoPost(url string, request, response interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "DoPushAppointmentRequest marshal")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "DoPushAppointmentRequest marshal")
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("appKey", app_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "DoPushAppointmentRequest post")
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return errors.Wrap(err, "DoPushAppointmentRequest status code is not 200")
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "DoPushAppointmentRequest readAll")
	}
	err = json.Unmarshal(data, response)
	if err != nil {
		return errors.Wrap(err, "DoPushAppointmentRequest unmarshal")
	}
	return nil
}
