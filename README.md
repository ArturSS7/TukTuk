
<p align="center">
 <img src="https://user-images.githubusercontent.com/19762721/91321586-74a99380-e7c7-11ea-8cb8-f5f25907a33a.png">
</p>

# TukTuk

This project was initially started as a part of [Digital Security](https://github.com/DSecurity)'s internship ["Summer of Hack 2020"](https://dsec.ru/about/traineeship/).

TukTuk is an open source tool that designed to make pentester's life easier by catching and logging different types of requests. TukTuk is written in Go, but has a little part of python code. 
## How to install
```TODO```

## Supported protocols:

 - HTTP
 - HTTPS
 - DNS
 - FTP
 - LDAP
 - SMTP (part of the code taken from [go-smtp](https://github.com/emersion/go-smtp))
 - SMB (used [impacket](https://github.com/SecureAuthCorp/impacket)'s SMB realisation)
 - Plain TCP
 
## Alerting 
TukTuk is designed not only to log requests, but to alert in case of it. Current alert options is:
 - Telegram
 - GMail
 
 Additional alert types can be added by writing a module. Feel free to make a pull request!

## Web interface
TukTuk is featuring a little web interface where user can manage some of the settings or look for logged request.
![Web interface](https://user-images.githubusercontent.com/19762721/91326276-d15b7d00-e7cc-11ea-8055-6760163c8dce.png)
