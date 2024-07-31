rm -rf nginx/keys
mkdir nginx/keys
cd nginx/keys
openssl genrsa -out server.key 4096
openssl req -out server.csr -key server.key -new
openssl x509 -req -days 3650 -signkey server.key -in server.csr -out server.crt