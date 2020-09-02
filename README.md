
<p align="center">
 <img src="https://user-images.githubusercontent.com/19762721/91321586-74a99380-e7c7-11ea-8cb8-f5f25907a33a.png">
</p>

# TukTuk

This project was initially started as a part of [Digital Security](https://github.com/DSecurity)'s internship ["Summer of Hack 2020"](https://dsec.ru/about/summerofhack/
).

TukTuk is an open source tool that designed to make pentester's life easier by catching and logging different types of requests. TukTuk is written in Go, but has a little part of python code. 

Also if you wondering why project naming is so strange - *TukTuk* stands for *Knock-Knock* in Russian.
## How to install

#### Requirements

- Golang >= 1.14.2
- PostgreSQL >= 12.1
- DNS domain
- If you want SMB protocol to work you will need to install impacket fork. 
You can do this with pip
```pip3 install git+https://github.com/ArturSS7/impacket.git@master```

#### Setting up DNS

In order to set up DNS do the following:
- Make A record ns.example.com which points your ip
- Make NS record e.example.com with contents of ns.example.com
- Make A record on example.com which points your ip
- Make MX record on example.com

Example shows how to setup TukTuk for *.h.example.com if your VPS' IP is 1.3.3.7. Main DNS is Cloudflare in our case, but you can use what you want
![DNS setup](https://user-images.githubusercontent.com/52138851/91592820-cbe66a00-e967-11ea-8a1f-e16379867ac6.png)

#### Building project and setting up database

Just run two scripts:
- install.sh
- db_init.sh

After doing the project will be moved to $GOPATH/src/TukTuk
You can run it with ./TukTuk

#### Configuring

The example configuration file is located in ```config/Config.json.example```.
There you can configure your domain settings, credentials and alerts
Move the file to ```Config.json``` if you are going to run the project.
Please change default credentials.

#### Configuring alerts

```TODO```

#### Getting HTTPS certificate

You will have to get a wildcard certificate for your domain.
You can do this with cert-bot.
A good article which will help you is [here](https://medium.com/@saurabh6790/generate-wildcard-ssl-certificate-using-lets-encrypt-certbot-273e432794d7).
First start the app and then start the bot.
During setting up certificates cert-bot will ask you to add TXT challenge to you domain. Add the TXT challenge in the ```Config.json``` file and continue.
After getting certificate put its path to the config file.

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
TukTuk is designed not only to log requests, but to alert in case of it. Current alert options are:
 - Telegram
 - GMail
 
 Additional alert types can be added by writing a module. Feel free to make a pull request!

## Web interface
TukTuk is featuring a little web interface where user can manage some of the settings or look for logged request.
![Web interface](https://user-images.githubusercontent.com/19762721/91326276-d15b7d00-e7cc-11ea-8055-6760163c8dce.png)
