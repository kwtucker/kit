# Token

Token is a tool to get access rights from the Comcast MPX Console. It will lists accounts that are avalible to you including the account number, name, and access role you have for each.

## Getting Started

Clone this repo and install
 
```shell
$ go get github.com/kwtucker/kit/cli/token
$ cd $GOPATH/src/github.com/kwtucker/kit/cli/token
$ go install
```

## Run

```shell
$ token
```

## Help

```shell
$ token -h

Usage of token:
  -d Get user accounts and roles.
  -n string
     Name of user. ENVS -> MPXUSERNAME_ + NAME, MPXPASSWORD_ + NAME
```