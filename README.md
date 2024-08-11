# vrc-social-notifier

## Build and Run

### Windows

```powershell
go build -o cli.exe cmd/cli/main.go
$env:VRC_AUTH_TOKEN = 'YOUR_TOKEN'
.\cli.exe
```

* Pass `VRC_AUTH_TOKEN` environment variable

## Features

* Notify changing friend's status
* Notify friends joining current instance
