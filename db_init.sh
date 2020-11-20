#!/bin/bash
echo "[*] Setting up database"
echo "If you are having problems with running this script do \"su postgres\" and try again"

psql -c "create database tuktuk;"
if test -z "$1"
then
	echo "Please provide db username as 1 arg"
	exit 1
fi
if test -z "$2"
then 
	echo "Please provide db pass as 2 arg"
	exit 1
fi

psql -c "create user $1 with encrypted password '$2';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE tuktuk TO $1;"
psql -c "ALTER USER $1 WITH SUPERUSER;"
psql -U $1 -d tuktuk -c "create table http(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table ftp(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table https(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table tcp(id serial, data varchar, source_ip varchar, time varchar, port int);"
psql -U $1 -d tuktuk -c "create table ldap(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table smtp(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table smb(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table dns(id serial, data varchar, source_ip varchar, time varchar);"
psql -U $1 -d tuktuk -c "create table dns_domains(id serial, domain varchar unique, delete_time bigint);"

echo "[*] If you have received no errors during database init everything is fine"
echo "Now you will have to set up https certificates for your domain"
