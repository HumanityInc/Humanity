package app

import (
	"../config"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/kurrik/oauth1a"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) wAuth() {

	if len(c.path) > 1 {

		switch c.path[1] {
		case "gplus":
			c.wGooglePlus()
		case "fb":
			c.wFacebook()
		case "tw":
			c.wTwitter()
		case "twcb":
			c.wTwitterCallback()
		}
	}
}

// ---

var (
	GP_CLIENT_ID       string
	GP_CLIENT_SECRET   string
	FB_APP_ID          string
	FB_APP_SECRET      string
	TW_CONSUMER_KEY    string
	TW_CONSUMER_SECRET string

	GP_REDIRECT_URI = "http://test.ishuman.me:2001/auth/gplus"
	FB_REDIRECT_URI = "http://test.ishuman.me:2001/auth/fb"
	TW_REDIRECT_URI = "http://test.ishuman.me:2001/auth/twcb"

	twitterConfig *oauth1a.Service
)

func init() {

	conf := config.GetConfig()

	TW_CONSUMER_KEY = conf.OAuth.Twitter.ConsumerKey
	TW_CONSUMER_SECRET = conf.OAuth.Twitter.ConsumerSecret

	FB_APP_ID = conf.OAuth.Facebook.AppId
	FB_APP_SECRET = conf.OAuth.Facebook.AppSecret

	GP_CLIENT_ID = conf.OAuth.Google.ClientId
	GP_CLIENT_SECRET = conf.OAuth.Google.ClientSecret

	twitterConfig = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    TW_CONSUMER_KEY,
			ConsumerSecret: TW_CONSUMER_SECRET,
			CallbackURL:    TW_REDIRECT_URI,
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
}

// ---

func (c *Client) wTwitterCallback() {

	requestTokenKey := c.req.FormValue("oauth_token")
	requestTokenSecret := ""

	if item, err := mc.Get("tw:" + requestTokenKey); err == nil {

		requestTokenSecret = string(item.Value)

		userConfig := &oauth1a.UserConfig{
			RequestTokenKey:    requestTokenKey,
			RequestTokenSecret: requestTokenSecret,
		}

		token, verifier, err := userConfig.ParseAuthorize(c.req, twitterConfig)

		if err != nil {
			http.Error(c.res, "Problem parsing authorization", 500)
			return
		}

		httpClient := new(http.Client)

		if err = userConfig.GetAccessToken(token, verifier, twitterConfig, httpClient); err != nil {
			http.Error(c.res, "Problem getting an access token", 500)
			return
		}

		fmt.Fprintf(c.res, "Access Token: %#v\n", userConfig)

		// fmt.Fprintf(c.res, "Access Token: %v\n", userConfig.AccessTokenKey)
		// fmt.Fprintf(c.res, "Token Secret: %v\n", userConfig.AccessTokenSecret)
		fmt.Fprintf(c.res, "Screen Name:  %v\n", userConfig.AccessValues.Get("screen_name"))
		fmt.Fprintf(c.res, "User ID:      %v\n", userConfig.AccessValues.Get("user_id"))

		return

	} else {

		// error
	}
}

func (c *Client) wTwitter() {

	httpClient := new(http.Client)
	userConfig := &oauth1a.UserConfig{}

	if err := userConfig.GetRequestToken(twitterConfig, httpClient); err != nil {
		http.Error(c.res, "Problem getting the request token", 500)
		return
	}

	url, err := userConfig.GetAuthorizeURL(twitterConfig)
	if err != nil {
		http.Error(c.res, "Problem getting the authorization URL", 500)
		return
	}

	mc.Set(&memcache.Item{
		Key:        "tw:" + userConfig.RequestTokenKey,
		Value:      []byte(userConfig.RequestTokenSecret),
		Expiration: 180,
	})

	// fmt.Printf("RequestTokenKey:    %#v\n", userConfig.RequestTokenKey) // public
	// fmt.Printf("RequestTokenSecret: %#v\n", userConfig.RequestTokenSecret)

	c.redirect(url)
}

// ---

func (c *Client) wGooglePlus() {

	conf := &oauth2.Config{
		ClientID:     GP_CLIENT_ID,
		ClientSecret: GP_CLIENT_SECRET,
		RedirectURL:  GP_REDIRECT_URI,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	code := c.req.FormValue("code")

	if code == "" {

		c.redirect(conf.AuthCodeURL("state"))
		return

	} else {

		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			return
		}

		response, err := conf.Client(oauth2.NoContext, tok).Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			fmt.Println("Get:", err)
			return
		}
		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)

		// fmt.Fprintf(c.res, "%s", data)

		// {
		// "id"
		// "email"
		// "verified_email"
		// "name"
		// "given_name"
		// "family_name"
		// "link"
		// "picture"
		// "gender"
		// }

		type Profile struct {
			UserId     string `json:"id"`
			Email      string `json:"email"`
			Name       string `json:"name"`
			GivenName  string `json:"given_name"`
			FamilyName string `json:"family_name"`
			Link       string `json:"link"`
			Picture    string `json:"picture"`
			Gender     string `json:"gender"`
		}

		profile := Profile{}

		err = json.Unmarshal(data, &profile)

		fmt.Fprintf(c.res, "%#v", profile)

		return
	}

	c.InternalServerError("")
}

// ---

func fbGetUrl(name, path string, param *url.Values) string {
	domainMap := map[string]string{
		"www":   "https://www.facebook.com/",
		"graph": "https://graph.facebook.com/",
		// "api":         "https://api.facebook.com/",
		// "api_read":    "https://api-read.facebook.com/",
		// "api_video":   "https://api-video.facebook.com/",
		// "graph_video": "https://graph-video.facebook.com/",
	}
	url, ok := domainMap[name]
	if !ok {
		return ""
	}
	url = url + path
	if param != nil {
		url = url + "?" + param.Encode()
	}
	return url
}

func randString(int) string {
	return fmt.Sprint(time.Now().UnixNano())
}

func (c *Client) wFacebook() {

	params := url.Values{}
	code := c.req.FormValue("code")
	error_param := c.req.FormValue("error")
	error_reason := c.req.FormValue("error_reason")

	if error_param == "" && code == "" {

		params.Add("client_id", FB_APP_ID)
		params.Add("redirect_uri", FB_REDIRECT_URI)
		params.Add("scope", "public_profile,email")
		params.Add("state", randString(14))

		url_get := fbGetUrl("www", "dialog/oauth", &params)

		// fmt.Println(url_get)

		c.redirect(url_get)
		return

	} else if code != "" {

		params.Add("client_id", FB_APP_ID)
		params.Add("client_secret", FB_APP_SECRET)
		params.Add("redirect_uri", FB_REDIRECT_URI)
		params.Add("code", code)

		url_get := fbGetUrl("graph", "oauth/access_token", &params)
		fmt.Println(url_get)

		response, err := http.Get(url_get)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)

		result, _ := url.ParseQuery(string(data))
		access_token := result.Get("access_token")
		// expires := result.Get("expires")

		if access_token != "" {

			params = url.Values{}
			params.Add("access_token", access_token)
			params.Add("fields", "id,name,first_name,last_name,email,locale,timezone,gender")
			url_get = fbGetUrl("graph", "me", &params)

			response, err := http.Get(url_get)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer response.Body.Close()

			data, err := ioutil.ReadAll(response.Body)

			type Profile struct {
				UserId     string `json:"id"`
				Email      string `json:"email"`
				Name       string `json:"name"`
				FirstName  string `json:"first_name"`
				SecondName string `json:"last_name"`
				Picture    string `json:"picture"`
				Gender     string `json:"gender"`
				Locale     string `json:"locale"`
				Link       string `json:"link"`
				Timezone   int    `json:"timezone"`
			}

			var profile Profile
			json.Unmarshal(data, &profile)

			/* {
				"id"
				"name"
				"first_name"
				"last_name"
				"email"
				"locale"
				"timezone"
				"gender"
			} */

			fmt.Fprintf(c.res, "%s\n", data)

			fmt.Fprintf(c.res, "%#v\n", profile)
		}

		return

	} else if error_param != "" {

		fmt.Fprintln(c.res, error_param)
		fmt.Fprintln(c.res, error_reason)
	}
}
