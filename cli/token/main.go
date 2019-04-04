package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	Detail bool
	Name   string
)

func main() {
	flag.BoolVar(&Detail, "d", false, "Get user accounts and roles.")
	flag.StringVar(&Name, "n", "", "Name of user. ENVS -> MPXUSERNAME_ + NAME, MPXPASSWORD_ + NAME")
	flag.Parse()

	envUsername := os.Getenv("MPXUSERNAME")
	envPassword := os.Getenv("MPXPASSWORD")

	if Name != "" {
		nameCaps := strings.ToUpper(Name)
		envUsername = os.Getenv("MPXUSERNAME" + "_" + nameCaps)
		envPassword = os.Getenv(fmt.Sprintf("MPXPASSWORD_%s", nameCaps))
	}

	if envUsername == "" || envPassword == "" {
		fmt.Println("Please make sure env variables MPXUSERNAME_NAME, MPXPASSWORD_NAME are set or MPXUSERNAME and MPXPASSWORD for use without passing in a name.")
		return
	}

	fmt.Printf(fmt.Sprintf("\nUsing %s: %s\n", strings.ToTitle(Name), envUsername))

	var ret string
	ret = fmt.Sprintf("%v:%v", envUsername, envPassword)
	ret = base64.StdEncoding.EncodeToString([]byte(ret))
	authheader := fmt.Sprintf("Basic %s", ret)

	signIn, err := mPXDataServiceLogin(authheader)
	if err != nil {
		fmt.Println(err)
	}

	if signIn != nil {
		fmt.Println(fmt.Sprintf("\nToken: %s", signIn.SignInResponse.Token))
		fmt.Println(fmt.Sprintf("User ID: %s", signIn.SignInResponse.UserID))

		if !Detail {
			return
		}

		byt, err := json.MarshalIndent(GetUserRolesForAccounts(signIn, GetUserAccountList(signIn.SignInResponse.Token)), "", "  ")
		if err != nil {
			fmt.Println("Could not marshal account list")
		}
		fmt.Println(fmt.Sprintf("\nUser Accounts: %s\n", string(byt)))
	}

}
