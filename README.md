# `service-check`
A quick and dirty web server that will generate a `200` response if the provided service is running, and a `5XX` if it is not.

## How To Use
Build `main.go` and then run `./service-check -svc "YOUR_SERVICE_NAME"`. You can also specify a custom port to serve your response with a `-p` flag. 

***EXAMPLE***: `./service-check -svc "thermald.service" -p "7777"`