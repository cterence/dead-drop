# dead drop.

my recreation of a dead drop application to securely share information. it's live on [dead-drop.terence.cloud](https://dead-drop.terence.cloud).

you can use it to share private informations to your peers. your data is encrypted/decrypted in the browser. the server cannot decrypt your data as it never receives the encryption key.

![Screencast from 2024-07-16 00-01-02](https://github.com/user-attachments/assets/534d9bed-e64d-45c6-8d70-042152c76534)

this project was made using:
- [go](https://go.dev/), [cobra](https://github.com/spf13/cobra) & [viper](https://github.com/spf13/viper/)
- [libsql-server](https://github.com/tursodatabase/libsql/tree/main/libsql-server)
- [sjcl](https://bitwiseshiftleft.github.io/sjcl/)
- [htmx](https://htmx.org/)
- [templ](https://github.com/a-h/templ)
- [tailwind](https://tailwindcss.com/)

try it out by running `docker compose --profile prod up` and go to http://localhost:3000.

## todo

- [x] generate a random password
- [x] store the encrypted message (libsql ? valkey ?)
- [x] generate a link to retrieve the message
- [x] create cool UI
- [x] ttl for drops, periodic flush
- [x] graceful errors when db is down
- [x] convert go code to cli
- [x] good (enough) logging
- [ ] ??????
- [ ] profit
