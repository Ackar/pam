# PAM, the Postgres And MySQL CLI

![pam-logo](https://user-images.githubusercontent.com/1406778/75760271-e3b7a980-5d9b-11ea-976d-1255d67448e9.png)

PAM is a simple command line tool to access Postgres and MySQL databases.

Features include:
- Auto-completion
- Coloration
- Smart display of wide results
- History
- Single binary with no dependencies

![demo](https://user-images.githubusercontent.com/1406778/75761479-064ac200-5d9e-11ea-8e20-565629adf4e2.gif)


### Installation

`GO111MODULE=on go get -u github.com/Ackar/pam`

### Configuration

Add your databases to `~/.pam.json` with the following format:

```json
{
  "mysqldb": {
    "type": "mysql",
    "dsn": "user:password@tcp(localhost:3306)/my_db"
  },
  "postgresdb": {
    "type": "postgres",
    "dsn": "postgres://user:password@localhost/my_db"
  }
}
```

You can then use `pam DB_NAME` to connect to any database.

#### Contributions welcome!
