package githubapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type OrganisationMember struct {
	Login     string `json:"login"`
	HtmlUrl   string `json:"html_url"`
	otherKeys map[string]any
}

type User struct {
	Name      string `json:"name"`
	HtmlUrl   string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
	otherKeys map[string]any
}

func GetUser(login string) (User, error) {
	path, err := url.JoinPath("https://api.github.com/users/", login)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot create path"), err)
	}

	var user User
	resp, err := http.Get(path)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot execute the GET request"), err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot read response"), err)
	}

	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return User{}, errors.Join(fmt.Errorf("Cannot read the response %s", string(bytes)), err)
	}

	return user, nil
}

const maxConcurrentRequests = 5

func GerOrganisationMembers(orgName string) ([]User, error) {
	path, err := url.JoinPath("https://api.github.com/orgs/", orgName, "/members")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create path"), err)
	}

	var members []OrganisationMember
	resp, err := http.Get(path)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute the GET request"), err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot read response"), err)
	}

	err = json.Unmarshal(bytes, &members)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot read the response %s", string(bytes)), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(maxConcurrentRequests)

	users := make([]User, len(members))
	for i, member := range members {
		wg.Add(1)
		go func() {
			sem.Acquire(ctx, 1)
			defer wg.Done()
			defer sem.Release(1)

			user, err := GetUser(member.Login)
			if err != nil {
				slog.Error("Could not read Github user - returning partial response", "err", err)
				users[i] = User{
					Name:    member.Login,
					HtmlUrl: member.HtmlUrl,
				}
			} else {
				users[i] = user
			}
		}()
	}

	wg.Wait()
	return users, nil
}
