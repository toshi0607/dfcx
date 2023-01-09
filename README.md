# dfcx

## Description

`dfcx` is a deployment tool of [Google Cloud Dialogflow CX](https://cloud.google.com/dialogflow/cx/docs) for multiple projects (dev, stg, and prd structure).

## Prerequisites

- Get permissions
  - export agent for base project
  - restore agent 
  - create version
  - update environment
- Set environment variables

```shell
export DF_LOCATION=asia-northeast1
export DF_BASE_PROJECT=xxx-dev
export DF_BASE_AGENT=111-222-333-444-555
export DF_BASE_ENV=222-333-444-555-666

export DF_STG_PROJECT=xxx-stg
export DF_STG_AGENT=333-444-555-666-777
export DF_STG_ENV=444-555-666-777-888

export DF_PRD_PROJECT=xxx-prd
export DF_PRD_AGENT=555-666-777-888-999
export DF_PRD_ENV=666-777-888-999-000
```

## Usage

```shell
# dev
$ dfcx agent deploy -v x.x.x base

# stg
$ dfcx agent deploy -v x.x.x stg 

# prd
$ dfcx agent deploy -v x.x.x prd
```

## Options

```shell
NAME:
   dfcx - operate dialogflow cx

USAGE:
   dfcx [global options] command [command options] [arguments...]

COMMANDS:
   agent    
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value, -l value  Log level [debug|info|warn|error] (default: "info") [$DF_LOG_LEVEL]
   --help, -h                   show help (default: false)

---

NAME:
   dfcx agent

USAGE:
   dfcx agent command [command options] [arguments...]

DESCRIPTION:
   dialogflow cx agent

COMMANDS:
   deploy   
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --location value      agent location (default: "asia-northeast1") [$DF_LOCATION]
   --base-project value  base project name [$DF_BASE_PROJECT]
   --base-agent value    base agent ID [$DF_BASE_AGENT]
   --base-env value      base environment ID [$DF_BASE_ENV]
   --help, -h            show help (default: false)

---

NAME:
   dfcx agent deploy

USAGE:
   dfcx agent deploy command [command options] [arguments...]

DESCRIPTION:
   Operate database

COMMANDS:
   base     
   stg      
   prd      
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --version value, -v value  version
   --help, -h                 show help (default: false)

---
```

## Installation

```shell
$ go install github.com/toshi0607/dfcx@latest
```

## License

[MIT](./LICENSE)
