# biter
A simple command tool used for campus network login in BIT 

![Golang](https://img.shields.io/github/go-mod/go-version/Nonsensersunny/biter)

![Travis](https://travis-ci.com/Nonsensersunny/biter.svg?branch=master)

## Installing

Using biter is easy. Several ways are provided to install the tool. 

1. Via compiled executable file

2. Compile by yourself

   ```shell
   # ensure that GOMODULE111MODULE has been set to `on`
   $ make prepare && make [linux|windows|darwin]
   ```

3. Run the code directly

   ```shell
   # make sure that dependencies are installed by `make prepare`
   $ go run cmd/main.go COMMANDS [OPTIONS]
   ```

## Commands

- Run `biter help` to see all commands contained. For example

  ```shell
  $ biter help
  biter V0.1.1
  A simple tool for login to campus network in BIT
  
  Usage:  biter COMMAND [OPTIONS]
  
  Options:
    -a        customize account configuration file path
    -c        customize account settings
    -g        customize full configuration file path
    -h        customize server configuration file path
  
  Commands:
    config    Account settings    
    login     Login to network    
    logout    Network logout      
    account   Get account info    
    update    Update biter tool  
  ```

## Getting started

Before using `biter` tool, account or server configurations need to be setup. Several ways are provided

- Command

  ```shell
  $ biter config
  Username:3220191000
  Password:123456789
  ```

- $ROOT/.biter/config.yaml (Create the file if not exists)

  ```yaml
  http:
    portal: 10.0.0.55
    challenge-url: cgi-bin/get_challenge
    srun-portal-url: cgi-bin/srun_portal
    succeed-url-origin: srun_portal_pc_succeed.php
    succeed-url-cmcc: srun_portal_pc_succeed_yys.php
    succeed-url-wcdma: srun_portal_pc_succeed_yys_cucc.php
  
  basic:
    username: 3220191000
    password: 123456789
  ```

- Specify configuration file

  ```shell
  $ biter -a [ACCOUNT_SETTINGS_FILE] -g [GLOBAL_SETTINGS_FILE] -s [SERVER_SETTINGS_FILE]
  ```