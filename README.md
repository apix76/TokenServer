#Description
Create and refresh tokens (Access, Refresh). Send via http/https in formate base64. Contains in Postgresql with bycrypt. If users ip is mismatch with ip from tokens, server send letter "email warning" to users mail. 

#How get 

For get code execute command below:
```bash
go clone github.com/apix76/TokenServer@latest <Target_Path>
```
##Build with Go build
For build code with go you can command below:
```bash
go build <Path-to-project>/TokenServer
```

##Docker
For create image from Dockerfile via Docker use this:
```bash
sudo docker build -t tokensserver <path-to-project>/TokenServer
```
for run image in the docker use command below:
```bash
sudo docker run --network=host -p <HttpPort>:<HttpPort> -p <HttpsPort>:<HttpsPort> --volume <Path-to-project>/TokenServer/config.cfg:/TokenServer/config.cfg tokensserver
```
How set config file, watch below ↓

##Request example
Send a requests body in json format. 

Create tokens
```bash
Request ↓
curl http(s)://<DomensServer>:<Http(s)Port>/<path_get> -d '{"guid":"<guid>"}'

Respons↓
{"RefreshToken":"<RefreshToken>","AccessToken":"<AccessToken>"}
```

Refresh tokens
```bash
Request ↓
curl http(s)://<DomensServer>:<Http(s)Port>/<path_refresh> -d '{"RefreshToken":"<RefreshToken>"}'

Respons↓
{"RefreshToken":"<RefreshToken>","AccessToken":"<AccessToken>"}
```

##Example set config.cfg
```json
{
  "CertFile":"<Path>/<CertFile>",
  "Keyfile":"<Path>/<Keyfile>",
  "HttpPort":":8080",
  "HttpsPort":":8081",
  "GetPath":"/get",
  "RefreshPath":"/refresh",
  "PgsqlNameServe":"postgres://<PsqlUserName>:<PsqlUserPassword>@<PsqlDomen>:<PsqlServerPort>/<PsqlDataBaseName>",
  "ExpTimeAccess": 30,
  "ExpTimeRefresh": 30,
  "MailHost":"smtp.mail.ru",
  "MailHostPortSmtp": 465,
  "MailUserName":"<UserMailName>",
  "MailPassword":"<UserMailPassword>"
}
```



