# ft-analyser-bot
FT Analyser Bot helps to automate the NPS analysis and FT weekly analysis process.

# FT Analyser Bot Infra
FT Analyser Bot service is listening on port 2112.

## Pre-requisites
To start the ft-analyser-bot service, pre-requisites are:

1. Mac (Preferred)
2. Amplitude, Bugsnag, Bork, slackbot credentials


## Configurations
The configurations for the service are set using `config.yaml`. Sample of `config.yaml` is present at [config-dev.yaml](config-dev.yaml)

This config.yaml should be placed at `/etc/envs/config.yaml`, it contains: 

```
# Amplitude API credentials 
1. Amplitude - apikey, secretkey

# Bork token
2. Bork token

# Bugsnag API credentials
3. Bugsnag - projectID, authToken

# slack bot credentials
4. slackbot - authToken, appToken, channelID
```

## Build ft-analyser-bot service

Clone the repository, navigate to the cloned repository and download the dependencies using `go mod download`. Before building, ensure the `config.yaml` is configured accordingly and placed at required location.

To build the ft-analyser-bot binary, use the below command, ft-analyser-bot binary built using make is placed in `bin` directory.

```sh
# Using make, prefered for linux OS.
make build

# Using go build and run.
sudo go run cmd/main.go
```

## Run ft-analyser-bot service

`ft-analyser-bot` can be run using binary on mac machine.

### Using binary
To run service through binary, follow the below command:
```sh
# Start service
sudo ./bin/ft-analyser-bot
```

Now the serivce is started and listening on port 2112. 
* Logs for service can be found at `/var/log/pf9-ft/ft-analyser-bot.log`


## `ft-analyser-bot` APIs
To interact with the service, ft-analyser-bot APIs are needed.

```sh
# To get weekly analysis
curl --request GET --url 'http://localhost:2112/weeklyanalysis'| jq .

## Sample response
{
  "Total_User_Signups": 23,
  "User_Who_Verified_Emails": 23,
  "PrepNode_Details": {
    "PrepNode_Attempts": {
      "New_Users": 3,
      "Existing_Users": 12
    },
    "PrepNode_Success": {
      "New_Users": 3,
      "Existing_Users": 9
    },
    "PrepNode_Errors": {
      "New_Users": 0,
      "Existing_Users": 3
    }
  },
  "Cluster_Creation_Attempts": {
    "New_Users": 2,
    "Existing_Users": 10
  }
}


# To get nps analysis of a user
curl --request GET --url 'http://localhost:9112/npsanalysis/<USERID>' | jq .
```