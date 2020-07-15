# hakq
A basic golang server/client for distributing shell commands to run over multiple systems.

# Note
This is not ready for production use. It's an absolute bare-bones implementation of what I'm trying to achieve, and it's probably full of bugs. I'm hoping that opening this up to the community will mean that the project can be worked on by many, making it more stable and feature-rich over time.

# Installation

At present, this is just two golang files, so you need to build them:

```
git clone https://github.com/hakluke/hakq
cd hakq
go build client.go
go build server.go
```

You should now have two binary files, `client` and `server`.

# Usage

First you will need to create a certificate:
```
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
```
Note that as this cert is self-signed, you will need to use the `--insecure` flag on the client. As an alternative, you can actually create a signed certificate and use this securely.

## BIG FAT SECURITY WARNING
Running these scripts on your machine in an insecure way (i.e. using the --insecure flag, or using a bad password) is the equivalent to providing someone RCE on your machine. Be careful

On the server:
```
./server --port 1234 --password <your-secure-password>
```

On the client:
```
./client --server <server-hostname>:1234 --password <your-secure-password>
```
Note: if your cert is self-signed, you will need to also add `--insecure`, which is a bad idea.

# Contributions
Please contribute to this project - it needs many features and bug fixes. Check the issues tab for some ideas!
Contributions should be made by forking this project, adding your code to the fork, and then doing a pull request.

# Found a bug?
- Firstly, make sure you're using the most up to date code. If not, update and see if it has been fixed.
- Secondly, make sure you're using the latest golang. The most common thing I see people do it install a LTS version of Ubuntu, and then install the "latest" golang from the apt repos, which isn't actually the latest. Currently the latest is 1.14.5.
- If you're sure there's a bug in the latest code, create a issue on this repo. Even better, if you can fix it, make a PR!

# TLS
You will need to generate a key pair so that your comms between the server/client are encrypted. In order to do so, on the server, create a key pair in the current directory. You can achieve this using the following command (or just use --insecure if you're feeling lucky):

```
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem
```
