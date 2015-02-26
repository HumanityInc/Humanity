package paypal

import (
	_ "../db"
	"../model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	USE_SANDBOX = true
)

var (
	logger *log.Logger
)

func init() {

	file, err := os.OpenFile(`paypal_ipn.log`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	logger = log.New(file, ``, log.Ldate|log.Ltime|log.Lshortfile)
}

func InstantPaymentNotification(c *model.Client) {

	c.Req.ParseForm()

	paypalUrl := "https://www.paypal.com/cgi-bin/webscr"

	if USE_SANDBOX {
		paypalUrl = "https://www.sandbox.paypal.com/cgi-bin/webscr"
	}

	param := url.Values{}
	param.Add("cmd", "_notify-validate")

	for key, value := range c.Req.Form {
		param.Add(key, value[0])
	}

	// payment_date 05:11:11 24 Feb 2015 PST
	// mc_handling1 1.67
	// item_number1 AK-1234
	// mc_gross1 9.34
	// residence_country US
	// payer_email buyer@paypalsandbox.com
	// verify_sign AFcWxV21C7fd0v3bYYYRCpSSRl31A-YLQvjGk1lxd2JkxsBNgGpf-o6.
	// payment_status Completed
	// payer_id TESTBUYERID01
	// business seller@paypalsandbox.com
	// mc_gross 12.34
	// payment_type instant
	// last_name Smith
	// txn_type cart
	// test_ipn 1
	// notify_version 2.4
	// address_zip 95131
	// address_name John Smith
	// receiver_id seller@paypalsandbox.com
	// address_country United States
	// mc_shipping 3.02
	// address_country_code US
	// address_state CA
	// payer_status unverified
	// address_status confirmed
	// item_name1 something
	// invoice abc1234
	// tax 2.02
	// mc_currency USD
	// first_name John
	// mc_handling 2.06
	// receiver_email seller@paypalsandbox.com
	// address_city San Jose
	// address_street 123, any street
	// txn_id 840783238
	// mc_fee 0.44
	// mc_shipping1 1.02
	// custom xyz123

	// fmt.Println(len(c.Req.Form))

	dataForm, _ := json.Marshal(&c.Req.Form)

	func(rawData, jsonForm string) {

		client := http.Client{
			Timeout: time.Duration(60 * time.Second),
		}

		post, _ := http.NewRequest("POST", paypalUrl, bytes.NewBufferString(rawData))
		post.Header.Add("Content-Length", strconv.Itoa(len(rawData)))
		post.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		if response, err := client.Do(post); err == nil {

			defer response.Body.Close()

			if data, err := ioutil.ReadAll(response.Body); err == nil {

				text := string(data)

				logger.Println(text, "|", jsonForm)

				switch text {
				case "VERIFIED":

				case "INVALID":

				default:

				}

				// TODO write to DB

				fmt.Println(text)

			} else {
				fmt.Println(err)
				logger.Println(err, "|", jsonForm)
			}

		} else {
			fmt.Println(err)
			logger.Println(err, "|", jsonForm)
		}

	}(param.Encode(), string(dataForm))

	c.Res.Write([]byte{})
}
