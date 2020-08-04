psql -U postgres -c "create database tuktuk;"
psql -U postgres -c "create user tuk with password ZuppaSecurePwd;"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE tuktuk TO tuk;"
psql -U tuk -c "create table http(id serial, data varchar, source_ip varchar, time varchar);"

