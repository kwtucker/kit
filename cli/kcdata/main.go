package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	arg := new(ARG)
	flag.BoolVar(&arg.Verbose, "v", false, "Verbose. To list full objects and data.")
	flag.BoolVar(&arg.Objects, "obj", false, "Objects. No decoded data.")
	flag.BoolVar(&arg.Data, "data", false, "Decoded data.")
	flag.StringVar(&arg.Name, "name", "", "Get by name.")
	flag.StringVar(&arg.Secret, "secret", "", "Secret example: -name NAME -secret 'key=val,key=val'")
	flag.StringVar(&arg.Delete, "delete", "", "Delete secret example: -delete name")
	flag.Parse()

	if len(os.Args) < 2 {
		arg.Data = true
	}

	if arg.Delete != "" {
		DeleteSecret(arg.Delete)
		return
	}

	if arg.Name == "" {
		object, data, _, err := HandleSecretObjects(arg)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error HandleSecretObjects: %+v", err))
			return
		}
		HandlePrint(arg, object, data)
		return
	}

	if arg.Name != "" {
		if arg.Secret != "" {
			err := StoreSecret(arg.Name, arg.Secret)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error StoreSecret: %+v", err))
				return
			}

		}

		object, secret, err := GetSecretObj(arg.Name)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error GetSecretObj: %+v", err))
			return
		}

		data, _, err := GetDataDecoded(secret.Data)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error GetDataDecoded: %+v", err))
			return
		}

		HandlePrint(arg, object, data)
	}
}
