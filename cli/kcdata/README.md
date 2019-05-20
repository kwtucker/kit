# KCDATA

KCDATA is a cli that will retrieve your kubernetes secrets from your current context and decode them so you can see exactly what they are.

## Getting Started

#### Get this repo and install

```shell
$ go get github.com/kwtucker/kit/cli/kcdata
$ cd $GOPATH/src/github.com/kwtucker/kit/cli/kcdata
$ go install
```

## Note
Make sure you are in the desired namespace and context.


## Run

### Get Full Objects

```shell
$ kcdata -obj
```

### Get Secret By Name

```shell
$ kcdata -name ExampleName
```

### Create/Update Secret

```shell
$ kcdata -name ExampleName -secret "Key=SecretText,Key=SecretText"
```

### Delete Secret
```shell
$ kcdata -delete ExampleName
```

## Help

```shell
$ kcdata -h

Usage of kcdata:
  -data
    	Decoded data.
  -delete string
    	Delete secret example: -delete name
  -name string
    	Get by name.
  -obj
    	Objects. No decoded data.
  -secret string
    	Secret example: -name NAME -secret 'key=val,key=val'
  -v	Verbose. To list full objects and data.
```