package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

var baseURL *url.URL

func init() {
	var err error
	if baseURL, err = url.Parse("https://api.twitch.tv/helix/"); err != nil {
		panic(err)
	}
}

type API struct {
	config Config
	client *http.Client
}

type Config struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func New(ctx context.Context, cfg Config) (api *API, err error) {
	api = &API{
		config: cfg,
	}

	oauth := clientcredentials.Config{
		ClientID:     api.config.ClientID,
		ClientSecret: api.config.ClientSecret,
		Scopes:       api.config.Scopes,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	api.client = oauth.Client(ctx)

	return
}

func (api *API) NewRequest(path string, method string, body io.Reader) (req *http.Request, err error) {
	var url *url.URL
	if url, err = baseURL.Parse(path); err != nil {
		return nil, err
	}

	req, err = http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", api.config.ClientID)

	return
}

func ZipQuery(args ...interface{}) (query *url.Values, err error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("len(args)%%2 != 0")
	}

	query = &url.Values{}

	for i := 0; i < len(args); i += 2 {
		var key string
		var ok bool

		if key, ok = args[i+0].(string); !ok {
			return nil, fmt.Errorf("arg %d is not a string", i+0)
		}

		switch arg := args[i+1].(type) {
		case int:
			query.Set(key, strconv.Itoa(arg))
		case string:
			query.Set(key, arg)
		case []string:
			for _, v := range arg {
				query.Add(key, v)
			}
		case fmt.Stringer:
			query.Set(key, arg.String())
		default:
			return nil, fmt.Errorf("arg %d is an invalid type", i+1)
		}
	}

	return query, nil
}

func StructToQuery(s interface{}) (query *url.Values, err error) {
	value := reflect.ValueOf(s)
	typ := value.Type()
	l := typ.NumField()
	zip := make([]interface{}, 0, l*2)
	for i := 0; i < l; i++ {
		if value.Field(i).IsZero() {
			continue
		}
		key := typ.Field(i).Tag.Get("query")
		if key == "" {
			key = typ.Field(i).Name
			log.Println(key)
		}
		zip = append(zip, key, value.Field(i).Interface())
	}
	return Z	ipQuery(zip...)
}

func (api *API) dataRequest(method, path string, in interface{}, out interface{}) (data DataResponse, err error) {
	req, err := api.NewRequest(path, method, nil)
	if err != nil {
		return
	}

	query, err := StructToQuery(in)
	if err != nil {
		return
	}

	req.URL.RawQuery = query.Encode()
	res, err := api.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("server returned %d", res.StatusCode)
		return
	}

	data.Data = out
	err = json.NewDecoder(res.Body).Decode(&data)

	return data, err
}

func (api *API) GetUsers(ur UsersRequest) (users []User, data DataResponse, err error) {
	data, err = api.dataRequest("GET", "users", ur, &users)
	return
}

func (api *API) GetStreams(sr StreamsRequest) (streams []Stream, data DataResponse, err error) {
	data, err = api.dataRequest("GET", "streams", sr, &streams)
	return
}

func (api *API) GetUsersFollows(ufr UsersFollowsRequest) (follows []UsersFollows, data DataResponse, err error) {
	data, err = api.dataRequest("GET", "users/follows", ufr, &follows)
	return
}
