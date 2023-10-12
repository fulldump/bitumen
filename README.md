# Bitumen

// todo: what is bitumen?

## Interfaces

Bitumen exposes filesystem in different ways, http, ftp, sftp.

### HTTP interfce

Default address: `0.0.0.0:8080`.

// todo: write this

### FTP interface

Default address: `0.0.0.0:2121`.

FTP is not a secure interface, it is implemented only for compatibility reasons.
It can be used with a restricted set of folders and/or operations (read only,
write only, read-write).

By default, no directories are exposed via FTP.

Please, use it with caution.

// todo: write this

### SFTP interface

Default address: `0.0.0.0:2222`.

SFTP exposes a filesystem via SSH+FTP subsystem. It is a secure way to transfer
files and easy to integrate with your linux filesystem with `sshfs` command.

A pair of RSA keys is required.

```shell
ssh-keygen -t rsa -f id_rsa
```

How to connect with sftp:

```shell
sftp -P 2222 -o IdentityFile=/path/to/YOUR_CLIENT_RSA 127.0.0.1
```

How to mount with sshfs:

```shell
sshfs -p 2222 -o IdentityFile=/path/to/YOUR_CLIENT_RSA user@127.0.0.1:/ /your/mount/point
```

How to generate client side certificates:

```shell
ssh-keygen -t ed25519 -f /path/to/YOUR_CLIENT_RSA
```