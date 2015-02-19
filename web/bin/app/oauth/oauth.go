package oauth

import (
	"../../config"
	"../../db"
	"../../model"
	"../../session"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/kurrik/oauth1a"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	fbAppId, fbAppSecret                        string
	gpClientId, gpClientSecret                  string
	twConsumerKey, twConsumerSecret             string
	gpRedirectUri, fbRedirectUri, twRedirectUri string

	twitterConfig *oauth1a.Service
	mc            *memcache.Client
)

func init() {

	conf := config.GetConfig()

	fbAppId = conf.OAuth.Facebook.AppId
	fbAppSecret = conf.OAuth.Facebook.AppSecret
	fbRedirectUri = conf.OAuth.Facebook.RedirectUrl

	gpClientId = conf.OAuth.Google.ClientId
	gpClientSecret = conf.OAuth.Google.ClientSecret
	gpRedirectUri = conf.OAuth.Google.RedirectUrl

	twConsumerKey = conf.OAuth.Twitter.ConsumerKey
	twConsumerSecret = conf.OAuth.Twitter.ConsumerSecret
	twRedirectUri = conf.OAuth.Twitter.RedirectUrl

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

	mc = memcache.New(conf.Storage.Memcache...)

	if err := mc.Set(&memcache.Item{Key: "test", Value: []byte("test")}); err != nil {
		fmt.Println("memcache", err)
		os.Exit(1)
	}
}

// ---

func AuthTwitterCallback(c *model.Client) {

	// TODO denied

	denied := c.Req.FormValue("denied")

	if denied != "" {
		c.Redirect("/")
		return
	}

	requestTokenKey := c.Req.FormValue("oauth_token")
	requestTokenSecret := ""

	if item, err := mc.Get("tw:" + requestTokenKey); err == nil {

		requestTokenSecret = string(item.Value)

		userConfig := oauth1a.NewAuthorizedConfig(requestTokenKey, requestTokenSecret)

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

		// fmt.Fprintf(c.Res, "Access Token: %#v\n", userConfig)
		// fmt.Fprintf(c.Res, "Access Token: %v\n", userConfig.AccessTokenKey)
		// fmt.Fprintf(c.Res, "Token Secret: %v\n", userConfig.AccessTokenSecret)
		// fmt.Fprintf(c.Res, "Screen Name:  %v\n", userConfig.AccessValues.Get("screen_name"))
		// fmt.Fprintf(c.Res, "User ID:      %v\n", userConfig.AccessValues.Get("user_id"))

		userId := userConfig.AccessValues.Get("user_id")
		firstName := userConfig.AccessValues.Get("screen_name")

		user, err := db.GetUserBySocialId(userId, model.SN_TWITTER)

		if err == nil {

			ukey := session.SetUserCookie(c.Res)
			session.SetUser(*user, ukey)

			if user.Email == "" {
				c.Redirect("/#!email")
			} else {

				if db.Invitee[strings.ToLower(user.Email)] == 1 {
					c.Redirect("/#!congratulations")
				} else {
					c.Redirect("/#!login")
				}
			}

		} else {

			socialProfile := model.SocialProfile{
				Id:        userId,
				SnId:      model.SN_TWITTER,
				FirstName: firstName,
				LastIp:    c.Ip(),
			}

			if user, err = db.RegisterSocialUser(socialProfile); err == nil {

				ukey := session.SetUserCookie(c.Res)
				session.SetUser(*user, ukey)
				c.Redirect("/#!email")

			} else {
				// error
			}
		}

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
			c.Redirect("/")
			fmt.Println(err)
			return
		}

		response, err := conf.Client(oauth2.NoContext, token).Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			c.Redirect("/")
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

				user, err := db.GetUserBySocialId(googleProfile.UserId, model.SN_GOOGLEPLUS)

				if err == nil {

					ukey := session.SetUserCookie(c.Res)
					session.SetUser(*user, ukey)

					if db.Invitee[strings.ToLower(googleProfile.Email)] == 1 {
						c.Redirect("/#!congratulations")
					} else {
						c.Redirect("/#!login")
					}

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

						ukey := session.SetUserCookie(c.Res)
						session.SetUser(*user, ukey)
						c.Redirect("/#!success")

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

	c.Redirect("/")
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
	state := c.Req.FormValue("state")
	error_param := c.Req.FormValue("error")
	error_reason := c.Req.FormValue("error_reason")

	if error_param == "" && code == "" {

		ukey := session.SetUserCookie(c.Res)

		params.Add("state", ukey)
		params.Add("client_id", fbAppId)
		params.Add("redirect_uri", fbRedirectUri)
		params.Add("scope", "public_profile,email")

		url_get := fbGetUrl("www", "dialog/oauth", &params)

		c.Redirect(url_get)
		return

	} else if code != "" {

		if state != session.GetSession(c.Req) {
			c.Redirect("/")
		}

		params.Add("code", code)
		params.Add("client_id", fbAppId)
		params.Add("redirect_uri", fbRedirectUri)
		params.Add("client_secret", fbAppSecret)

		url_get := fbGetUrl("graph", "oauth/access_token", &params)
		// fmt.Println(url_get)

		response, err := http.Get(url_get)
		if err != nil {
			fmt.Println(err)
			c.Redirect("/")
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

			type FacebookProfile struct {
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

			facebookProfile := FacebookProfile{}

			// json.Unmarshal(data, &profile)
			// fmt.Fprintf(c.Res, "%#v\n", profile)

			if err = json.Unmarshal(data, &facebookProfile); err == nil {

				user, err := db.GetUserBySocialId(facebookProfile.UserId, model.SN_FACEBOOK)

				if err == nil {

					ukey := session.SetUserCookie(c.Res)
					session.SetUser(*user, ukey)

					if db.Invitee[strings.ToLower(facebookProfile.Email)] == 1 {
						c.Redirect("/#!congratulations")
					} else {
						c.Redirect("/#!login")
					}

				} else {

					socialProfile := model.SocialProfile{
						Id:        facebookProfile.UserId,
						SnId:      model.SN_FACEBOOK,
						Email:     facebookProfile.Email,
						FirstName: facebookProfile.FirstName,
						LastName:  facebookProfile.SecondName,
						Picture:   facebookProfile.Picture,
						Link:      facebookProfile.Link,
						Gender:    facebookProfile.Gender,
						LastIp:    c.Ip(),
					}

					if user, err = db.RegisterSocialUser(socialProfile); err == nil {

						ukey := session.SetUserCookie(c.Res)
						session.SetUser(*user, ukey)
						c.Redirect("/#!success")

					} else {
						// error
					}
				}

			} else {
				// error
			}
		}

		return

	} else if error_param != "" {

		c.Redirect("/")
		fmt.Println(error_param)
		fmt.Println(error_reason)
	}
}
