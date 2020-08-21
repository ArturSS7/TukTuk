psql -U postgres -c "create database tuktuk;"
psql -U postgres -c "create user tuk with encrypted password ZuppaSecurePwd;"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE tuktuk TO tuk;"
psql -U tuk -c "create table http(id serial, data varchar, source_ip varchar, time varchar);"
psql -U tuk -c "create table ftp(id serial, data varchar, source_ip varchar, time varchar);"
psql -U tuk -c "create table https(id serial, data varchar, source_ip varchar, time varchar);"
psql -U tuk -c "create table tcp(id serial, data varchar, source_ip varchar, time varchar, port int);"

