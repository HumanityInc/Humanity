package sendmail

import (
	"../config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

const (
	TEMPLATE_PATH = `../template/en/_email/`
	MANDRILL_API  = `https://mandrillapp.com/api/1.0`
	FROM_EMAIL    = `welcome@ishuman.me`
	FROM_NAME     = `Humanity`
)

type (
	MailActivate struct {
		Link string
	}

	MailSuccess struct {
		Link string
	}

	MandrillReturn struct {
		Email        string `json:"email"`
		Status       string `json:"status"`
		Id           string `json:"_id"`
		RejectReason string `json:"reject_reason"`
	}
)

var (
	logger      *log.Logger
	mandrillKey = ""
)

func init() {

	conf := config.GetConfig()

	mandrillKey = conf.Sendmail.Mandrill.ApiKey

	file, err := os.OpenFile(`sendmail_errors.log`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	logger = log.New(file, ``, log.Ldate|log.Ltime|log.Lshortfile)
}

//

func (m *MailActivate) Send(email string) {

	err := send("activate.txt", email, "EmailActivate", m)

	if err != nil {
		logger.Println(err)
	}
}

func (m *MailSuccess) Send(email string) {

	err := send("success.txt", email, "Welcome to Humanity! :)", m)

	if err != nil {
		logger.Println(err)
	}
}

//

func send(tmpl_file, email, subject string, data interface{}) (ret error) {

	tmpl, err := template.ParseFiles(TEMPLATE_PATH + tmpl_file)
	if err == nil {

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, data)

		if err == nil {

			text := buf.String()

			messages := map[string]interface{}{
				"key": mandrillKey,
				"message": map[string]interface{}{
					// "text":       text,
					"html":       text,
					"subject":    subject,
					"from_email": FROM_EMAIL,
					"from_name":  FROM_NAME,
					"to": []map[string]interface{}{
						map[string]interface{}{
							"email": email,
						},
					},
				},
			}

			raw_body, err := json.Marshal(messages)
			if err != nil {
				logger.Println(err)
				return
			}

			string_body := string(raw_body)

			client := http.Client{}

			request, _ := http.NewRequest(`POST`, MANDRILL_API+`/messages/send.json`, bytes.NewBufferString(string_body))
			request.Header.Add("Content-Length", strconv.Itoa(len(string_body)))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(request)
			if err != nil {
				logger.Println(err)
				return
			}
			defer response.Body.Close()

			data, _ := ioutil.ReadAll(response.Body)

			request_result := []MandrillReturn{}

			err = json.Unmarshal(data, &request_result)

			fmt.Println(err)
			fmt.Println(string(data))
			fmt.Println(request_result)

			// TODO processed results

			return

		} else {

			ret = err
		}
	}

	ret = err

	return
}
