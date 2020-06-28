package twitch

import (
	"context"
	"fmt"
	"os"
	"testing"
)

var api *API

func TestAPI(t *testing.T) {

	var err error

	api, err = New(context.Background(), Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		Scopes:       []string{"channel:read:stream_key"},
	})

	if err != nil {
		t.Fatal(err)
	}
}

var user User

func TestGetUsers(t *testing.T) {
	users, data, err := api.GetUsers(UsersRequest{
		Login: os.Getenv("TWITCH_LOGIN"),
	})

	t.Logf("%+v", data)

	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatal(fmt.Errorf("invalid number of users returned: %d, expected %d", len(users), 1))
	}

	user = users[0]
}

var follows []UsersFollows

func TestUsersFollows(t *testing.T) {

	req := UsersFollowsRequest{
		FromID: user.ID,
	}

	for {
		f, data, err := api.GetUsersFollows(req)

		if err != nil {
			t.Fatal(err)
		}

		follows = append(follows, f...)

		if data.Pagination.Cursor == "" {
			break
		} else {
			req.After = data.Pagination.Cursor
		}
	}
}

func TestStreams(t *testing.T) {
	logins := make([]string, len(follows))

	for i := range follows {
		logins[i] = follows[i].ToName
	}

	if len(logins) < 1 {
		// if you're here that means twitch will find ALL live streams
		// we shouldn't do that
		t.Fatal("no logins to search streams for")
	}

	t.Log(logins)

	req := StreamsRequest{
		UserLogin: logins,
	}

	var streams []Stream

	for {
		s, data, err := api.GetStreams(req)

		if err != nil {
			t.Fatal(err)
		}

		streams = append(streams, s...)

		if data.Pagination.Cursor == "" {
			break
		} else {
			req.After = data.Pagination.Cursor
		}
	}

	t.Log(streams)
}
