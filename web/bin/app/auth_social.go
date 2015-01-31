package app

import (
	"../config"
	"../db"
	"../model"
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

func Auth(c *model.Client) {

	if len(c.Path) > 1 {

		switch c.Path[1] {
		case "gplus":
			AuthGooglePlus(c)
		case "fb":
			AuthFacebook(c)
		case "tw":
			AuthTwitter(c)
		case "twcb":
			AuthTwitterCallback(c)
		}
	}
}

// ---

var (
	fbAppId, fbAppSecret            string
	gpClientId, gpClientSecret      string
	twConsumerKey, twConsumerSecret string

	gpRedirectUri = "http://test.ishuman.me:2001/auth/gplus"
	fbRedirectUri = "http://test.ishuman.me:2001/auth/fb"
	twRedirectUri = "http://test.ishuman.me:2001/auth/twcb"

	twitterConfig *oauth1a.Service
)

func init() {

	conf := config.GetConfig()

	fbAppId = conf.OAuth.Facebook.AppId
	fbAppSecret = conf.OAuth.Facebook.AppSecret

	gpClientId = conf.OAuth.Google.ClientId
	gpClientSecret = conf.OAuth.Google.ClientSecret

	twConsumerKey = conf.OAuth.Twitter.ConsumerKey
	twConsumerSecret = conf.OAuth.Twitter.ConsumerSecret

	twitterConfig = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    twConsumerKey,
			ConsumerSecret: twConsumerSecret,
			CallbackURL:    twRedirectUri,
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
}

// ---

func AuthTwitterCallback(c *model.Client) {

	requestTokenKey := c.Req.FormValue("oauth_token")
	requestTokenSecret := ""

	if item, err := mc.Get("tw:" + requestTokenKey); err == nil {

		requestTokenSecret = string(item.Value)

		userConfig := &oauth1a.UserConfig{
			RequestTokenKey:    requestTokenKey,
			RequestTokenSecret: requestTokenSecret,
		}

		token, verifier, err := userConfig.ParseAuthorize(c.Req, twitterConfig)

		if err != nil {
			http.Error(c.Res, "Problem parsing authorization", 500)
			return
		}

		httpClient := new(http.Client)

		if err = userConfig.GetAccessToken(token, verifier, twitterConfig, httpClient); err != nil {
			http.Error(c.Res, "Problem getting an access token", 500)
			return
		}

		fmt.Fprintf(c.Res, "Access Token: %#v\n", userConfig)

		// fmt.Fprintf(c.Res, "Access Token: %v\n", userConfig.AccessTokenKey)
		// fmt.Fprintf(c.Res, "Token Secret: %v\n", userConfig.AccessTokenSecret)
		fmt.Fprintf(c.Res, "Screen Name:  %v\n", userConfig.AccessValues.Get("screen_name"))
		fmt.Fprintf(c.Res, "User ID:      %v\n", userConfig.AccessValues.Get("user_id"))

		return

	} else {

		// error
	}
}

func AuthTwitter(c *model.Client) {

	httpClient := new(http.Client)
	userConfig := &oauth1a.UserConfig{}

	if err := userConfig.GetRequestToken(twitterConfig, httpClient); err != nil {
		http.Error(c.Res, "Problem getting the request token", 500)
		return
	}

	url, err := userConfig.GetAuthorizeURL(twitterConfig)
	if err != nil {
		http.Error(c.Res, "Problem getting the authorization URL", 500)
		return
	}

	mc.Set(&memcache.Item{
		Key:        "tw:" + userConfig.RequestTokenKey,
		Value:      []byte(userConfig.RequestTokenSecret),
		Expiration: 180,
	})

	// fmt.Printf("RequestTokenKey:    %#v\n", userConfig.RequestTokenKey) // public
	// fmt.Printf("RequestTokenSecret: %#v\n", userConfig.RequestTokenSecret)

	c.Redirect(url)
}

// ---

func AuthGooglePlus(c *model.Client) {

	conf := &oauth2.Config{
		ClientID:     gpClientId,
		ClientSecret: gpClientSecret,
		RedirectURL:  gpRedirectUri,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	code := c.Req.FormValue("code")

	if code == "" {

		c.Redirect(conf.AuthCodeURL("state"))
		return

	} else {

		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			return
		}

		response, err := conf.Client(oauth2.NoContext, token).Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			fmt.Println("Get:", err)
			return
		}
		defer response.Body.Close()

		if data, err := ioutil.ReadAll(response.Body); err == nil {

			type GoogleProfile struct {
				UserId     string `json:"id"`
				Email      string `json:"email"`
				Name       string `json:"name"`
				GivenName  string `json:"given_name"`
				FamilyName string `json:"family_name"`
				Link       string `json:"link"`
				Picture    string `json:"picture"`
				Gender     string `json:"gender"`
			}

			googleProfile := GoogleProfile{}

			if err = json.Unmarshal(data, &googleProfile); err == nil {

				// fmt.Fprintf(c.Res, "%#v", googleProfile)

				user, err := db.GetUserBySocialId(googleProfile.UserId, model.SN_GOOGLEPLUS)

				fmt.Println("> ", user, err)

				if err == nil {

					c.WriteJson(user)

				} else {

					socialProfile := model.SocialProfile{
						Id:        googleProfile.UserId,
						SnId:      model.SN_GOOGLEPLUS,
						Email:     googleProfile.Email,
						FirstName: googleProfile.GivenName,
						LastName:  googleProfile.FamilyName,
						Picture:   googleProfile.Picture,
						Link:      googleProfile.Link,
						Gender:    googleProfile.Gender,
						LastIp:    c.Ip(),
					}

					if user, err = db.RegisterSocialUser(socialProfile); err == nil {

						c.WriteJson(user)

					} else {
						// error
					}
				}

			} else {
				// error
			}

		} else {
			// error
		}

		return
	}

	InternalServerError(c)
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

func AuthFacebook(c *model.Client) {

	params := url.Values{}
	code := c.Req.FormValue("code")
	error_param := c.Req.FormValue("error")
	error_reason := c.Req.FormValue("error_reason")

	if error_param == "" && code == "" {

		params.Add("client_id", fbAppId)
		params.Add("redirect_uri", fbRedirectUri)
		params.Add("scope", "public_profile,email")
		params.Add("state", randString(14))

		url_get := fbGetUrl("www", "dialog/oauth", &params)

		// fmt.Println(url_get)

		c.Redirect(url_get)
		return

	} else if code != "" {

		params.Add("client_id", fbAppId)
		params.Add("client_secret", fbAppSecret)
		params.Add("redirect_uri", fbRedirectUri)
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
				"id":"785049524848954",
				"name":"\u0410\u043b\u0435\u043a\u0441\u0435\u0439 \u0410\u0447\u043a\u0430\u0441\u043e\u0432",
				"first_name":"\u0410\u043b\u0435\u043a\u0441\u0435\u0439",
				"last_name":"\u0410\u0447\u043a\u0430\u0441\u043e\u0432",
				"email":"al_ghost\u0040inbox.ru",
				"locale":"ru_RU",
				"timezone":2,
				"gender":"male"
			} */

			fmt.Fprintf(c.Res, "%s\n", data)

			fmt.Fprintf(c.Res, "%#v\n", profile)

			// if profile.UserId != "" {

			// u, err := users.GetByFb(profile.UserId)

			// fmt.Println(u, err)

			// if err != nil {

			// 	user := database.User{}

			// 	user.Name = profile.Name
			// 	user.Email = profile.Email
			// 	user.FbId = profile.UserId
			// 	user.Photo = "https://graph.facebook.com/" + profile.UserId + "/picture?type=square&height=200&width=200"
			// 	user.Login = fmt.Sprint("fb_", profile.UserId)

			// 	users.AddUser(&user)

			// 	session.Values["id"] = user.Id
			// 	session.Values["user"] = user

			// } else {

			// 	session.Values["id"] = u.Id
			// 	session.Values["user"] = *u
			// }
			// }

			// c.Redirect("/")
		}

		return

	} else if error_param != "" {

		fmt.Fprintln(c.Res, error_param)
		fmt.Fprintln(c.Res, error_reason)
	}
}
