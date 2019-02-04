package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func HandleSecretObjects(arg *ARG) ([]byte, []byte, *Secrets, error) {
	_, secrets, err := GetSecretObjs()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error GetSecretObj: %+v", err))
		return nil, nil, nil, err
	}

	d := []map[string]interface{}{}

	for _, val := range secrets.Items {
		innerd := map[string]interface{}{}
		innerd["name"] = val.Metadata.Name

		_, innerd["data"], err = GetDataDecoded(val.Data)
		if err != nil {
			return nil, nil, nil, err
		}
		d = append(d, innerd)
	}

	dataOut, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, nil, nil, err
	}

	object, err := json.MarshalIndent(secrets, "", "  ")

	return object, dataOut, secrets, err
}

func GetSecretObj(name string) ([]byte, *Item, error) {

	cmd := exec.Command("kubectl", "get", "secret", name, "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}
	secret := &Item{}
	err = json.Unmarshal(out, secret)
	return out, secret, err
}

func GetSecretObjs() ([]byte, *Secrets, error) {

	cmd := exec.Command("kubectl", "get", "secrets", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}
	secrets := &Secrets{}
	err = json.Unmarshal(out, secrets)
	return out, secrets, err
}

func GetDataDecoded(data map[string]string) ([]byte, map[string]string, error) {
	m := map[string]string{}
	for key, val := range data {
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, nil, err
		}

		m[key] = string(decoded)
	}

	out, err := json.MarshalIndent(m, "", "  ")
	return out, m, err
}

func StoreSecret(name string, value string) error {
	fmt.Println("name ", name)
	fmt.Println("value ", value)

	cmdDelete := exec.Command("kubectl", "delete", "secret", name)
	cmdDelete.Run()

	cmd := exec.Command("kubectl", "create", "secret", "generic", name)

	values := strings.Split(value, ",")
	fmt.Println("values ", values)

	for _, v := range values {
		cmd.Args = append(cmd.Args, "--from-literal="+v)
	}
	return cmd.Run()
}

func DeleteSecret(name string) error {
	cmdDelete := exec.Command("kubectl", "delete", "secret", name)
	return cmdDelete.Run()
}

func HandlePrint(arg *ARG, obj, data []byte) {
	if arg.Verbose {
		fmt.Println(fmt.Sprintf("%s", string(obj)))
		fmt.Println(fmt.Sprintf("Decoded Data: %s", string(data)))
		return
	}

	if arg.Name != "" {
		if !arg.Objects && !arg.Data {
			fmt.Println(fmt.Sprintf("%s", string(obj)))
			fmt.Println(fmt.Sprintf("Decoded Data: %s", string(data)))
			return
		}

		if arg.Objects {
			fmt.Println(fmt.Sprintf("%s", string(obj)))
		}

		if arg.Data {
			fmt.Println(fmt.Sprintf("Decoded Data: %s", string(data)))
		}
		return
	}

	if arg.Objects {
		fmt.Println(fmt.Sprintf("%s", string(obj)))
	}

	if arg.Data {
		fmt.Println(fmt.Sprintf("Decoded Data: %s", string(data)))
	}
}
