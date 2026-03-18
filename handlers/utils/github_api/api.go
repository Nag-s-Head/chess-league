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
	Login     string `json:"login"`
	Name      string `json:"name"`
	HtmlUrl   string `json:"html_url"`
	AvatarUrl string `json:"avatar_url"`
	otherKeys map[string]any
}

func GetUser(login string, apiKey string) (User, error) {
	path, err := url.JoinPath("https://api.github.com/users/", login)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot create path"), err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot create the request"), err)
	}

	if apiKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot execute the GET request"), err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot read response"), err)
	}

	var user User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return User{}, errors.Join(fmt.Errorf("Cannot read the response %s", string(bytes)), err)
	}

	return user, nil
}

const maxConcurrentRequests = 5

func GerOrganisationMembers(orgName string, apiKey string) ([]User, error) {
	path, err := url.JoinPath("https://api.github.com/orgs/", orgName, "/members")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create path"), err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create the request"), err)
	}

	if apiKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute the GET request"), err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot read response"), err)
	}

	var members []OrganisationMember
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
		go func(i int, member OrganisationMember) {
			sem.Acquire(ctx, 1)
			defer wg.Done()
			defer sem.Release(1)

			user, err := GetUser(member.Login, apiKey)
			if err != nil {
				slog.Error("Could not read Github user - returning partial response", "err", err)
				users[i] = User{
					Login:   member.Login,
					Name:    member.Login,
					HtmlUrl: member.HtmlUrl,
				}
			} else {
				if user.Name == "" {
					user.Name = member.Login
				}
				users[i] = user
			}
		}(i, member)
	}

	wg.Wait()
	return users, nil
}

func GetAuthenticatedUser(accessToken string) (User, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot create the request"), err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot execute the GET request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("Github API returned status %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, errors.Join(errors.New("Cannot read response"), err)
	}

	var user User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return User{}, errors.Join(fmt.Errorf("Cannot read the response %s", string(bytes)), err)
	}

	return user, nil
}

func IsMemberOfOrg(orgName string, login string, apiKey string) (bool, error) {
	path, err := url.JoinPath("https://api.github.com/orgs/", orgName, "/members/", login)
	if err != nil {
		return false, errors.Join(errors.New("Cannot create path"), err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return false, errors.Join(errors.New("Cannot create the request"), err)
	}

	if apiKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Join(errors.New("Cannot execute the GET request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("Unexpected status code from Github API: %d", resp.StatusCode)
}
