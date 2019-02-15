package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetUserRolesForAccounts(user *signInMPX, accountMap map[string]string) map[string]string {
	for account, accountname := range accountMap {
		userPermissions := GetUserPermissions(user, account)
		if userPermissions != nil && len(userPermissions.Entries) > 0 {

			permission := userPermissions.Entries[0]

			if len(permission.RoleIDs) > 0 {
				for _, roleid := range permission.RoleIDs {
					role := GetRole(roleid, user.SignInResponse.Token)
					if role != nil {
						accountMap[account] = "[ " + role.Title + " ] : " + accountname
					}
				}
			}

		} else {
			accountMap[account] = "[ NIL ] : " + accountname
		}
	}
	return accountMap
}

func GetUserPermissions(s *signInMPX, account string) *responseWrapperPerermission {
	permissionURL, err := url.Parse(getPermission)
	if err != nil {
		fmt.Println(err)
	}

	q := permissionURL.Query()

	q.Add("token", s.SignInResponse.Token)
	q.Add("byUserId", s.SignInResponse.UserID)
	q.Add("byOwnerId", account)

	permissionURL.RawQuery = q.Encode()

	resp, err := http.Get(permissionURL.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	getPermission := &responseWrapperPerermission{}
	err = json.Unmarshal(byt, &getPermission)

	return getPermission
}

func GetRole(id string, token string) *roleMPX {
	id = strings.Split(id, "Role/")[1]
	roleURL, err := url.Parse(getMpxRole)
	if err != nil {
		fmt.Println(err)
	}
	q := roleURL.Query()

	q.Add("token", token)
	q.Add("byId", id)

	roleURL.RawQuery = q.Encode()
	resp, err := http.Get(roleURL.String())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	getRole := &responseWrapperRole{}
	err = json.Unmarshal(byt, &getRole)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(getRole.Entries) < 1 {
		fmt.Println("no roles")
		return nil
	}

	return getRole.Entries[0]
}

// mPXDataServiceLogin Logs into MPX for a token, username, Id, duration, idletime.
func mPXDataServiceLogin(auth string) (*signInMPX, error) {
	var httpClient http.Client
	var request *http.Request
	var response *http.Response
	var contentBuffer io.Reader
	var err error

	request, err = http.NewRequest(http.MethodGet, getToken, contentBuffer)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to create request: %v\n", err)
		fmt.Println(errorMsg)
	}

	request.Header.Set("Authorization", auth)

	// Make GET request
	response, err = httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	byt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Fill the GetSelfMPX Stucture with response. If the response is empty fill error struct.
	signInResponse := &signInMPX{}
	err = json.Unmarshal(byt, &signInResponse)
	if err != nil || signInResponse.SignInResponse.Token == "" {

		// Returns a error response object and a error.
		// Error response object is logged in this function,
		// so it is generally ignored unless passing the error object somewhere else.
		_, logErr := logErrorMPXResponse(byt)
		if logErr != nil {
			fmt.Println(logErr)
		}

		if err != nil {
			fmt.Println(err)
		}

		return nil, err
	}
	return signInResponse, nil
}

func GetUserAccountList(token string) map[string]string {
	list := map[string]string{}

	byt, _, err := sendGet(getAccount + "&fields=title,id&token=" + token)
	if err != nil {
		fmt.Println(fmt.Sprintf("%+v", err))
	}

	type any struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
	}

	type accountlist struct {
		StartIndex   int   `json:"startIndex,omitempty"`
		ItemsPerPage int   `json:"itemsPerPage,omitempty"`
		EntryCount   int   `json:"entryCount,omitempty"`
		Entries      []any `json:"entries,omitempty"`
	}
	alist := &accountlist{}

	err = json.Unmarshal(byt, alist)
	if err != nil {
		fmt.Println(err)
	}

	if len(alist.Entries) > 0 {
		for _, val := range alist.Entries {
			acountIDIndex := strings.LastIndex(val.ID, "/")
			accountID := val.ID[acountIDIndex+1:]
			// list[val.Title] = accountID
			list[accountID] = val.Title

		}
		return list
	}
	return nil
}

// encodeCredentials - Encodes the username and password for GET request.
func encodeCredentials(user, auth string) string {
	var ret string
	ret = fmt.Sprintf("%v:%v", user, auth)
	ret = base64.StdEncoding.EncodeToString([]byte(ret))
	return fmt.Sprintf("Basic %s", ret)
}

// Send basic get request with token
func sendGet(url string) ([]byte, http.Header, error) {
	var httpClient http.Client
	var request *http.Request
	var response *http.Response
	var contentBuffer io.Reader
	var err error

	request, err = http.NewRequest(http.MethodGet, url, contentBuffer)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to create request: %v\n", err)
		return nil, nil, errors.New(errorMsg)
	}

	// Make GET request
	response, err = httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	byt, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return byt, response.Header, nil
}

func badArgLog(name string, accounts map[string]user) {
	fmt.Println(" ")
	if name != "" {
		fmt.Println(fmt.Sprintf("Not available: \"%s\"\n", name))
	}

	byt, _ := json.MarshalIndent(accounts, "", "  ")
	fmt.Println(fmt.Sprintf("Configured Accounts Include: %s", string(byt)))
	availableAccounts := []string{}
	for val := range accounts {
		availableAccounts = append(availableAccounts, val)
	}
	fmt.Println(fmt.Sprintf("\nRetry with one of these: \n- %s\n", strings.Join(availableAccounts, "\n- ")))
	fmt.Println("Example: mpx name\n")
	return
}

// logErrorMPXResponse Logs the
func logErrorMPXResponse(byt []byte) (*errorMPX, error) {
	errResponse := &errorMPX{}
	err := json.Unmarshal(byt, &errResponse)
	if errResponse.ResponseCode != 0 {
		fmt.Println("")
		fmt.Printf("Response Code:  %d\n", errResponse.ResponseCode)
		fmt.Printf("Description:    %s\n", errResponse.Description)
		fmt.Printf("Title:          %s\n", errResponse.Title)
		fmt.Printf("Correlation ID: %s\n", errResponse.CorrelationID)
		fmt.Printf("IsException:    %v\n", errResponse.IsException)
		fmt.Println("")
	}
	return errResponse, err
}
