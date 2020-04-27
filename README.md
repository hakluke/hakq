# hakq
A basic golang server/client for distributing tasks over multiple systems.

# TLS
You will need to generate a key pair so that your comms between the server/client are encrypted. In order to do so, on the server, create a key pair in the current directory. You can achieve this using the following command:

```
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
```
