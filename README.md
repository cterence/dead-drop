# dead drop

my recreation of a dead drop application to securely share information.

![Screencast from 2024-07-16 00-01-02](https://github.com/user-attachments/assets/534d9bed-e64d-45c6-8d70-042152c76534)

## todo

- [x] generate a random password
- [x] store the encrypted message (libsql ? valkey ?)
- [x] generate a link to retrieve the message
- [x] create cool UI
- [ ] good logging
- [ ] ttl for drops, periodic flush
- [ ] graceful errors when db is down
- [ ] convert go code to cli
