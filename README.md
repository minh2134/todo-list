# Simple TODO app

## How to run
Requirements:
- Go 1.26

Instructions: At root folder simply type the command:
```bash
$ CGO_ENABLED='1' go run .
```

And the application should be available at `localhost:3000`

## Config
You can specify where the DB file will stay by setting the environment variable `TODO_DB_FILE` to the **filepath** (not dir path) you wish to save. The default is the root folder.

## Docker
A Dockerfile is available at the root folder
