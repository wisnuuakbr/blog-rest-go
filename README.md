# Building Simple Blog App with Golang
simple RestAPI with Go, Gin, Gorm, and PostgreSQL

## Requirements

This simple apps is currently extended with the following requirements.  
Instructions on how to use them in your own application are linked below.

| Requirement | Version |
| ----------- | ------- |
| Go          | 1.21.5  |
| PostgreSQL  | 14.10.x |

## Installation

Make sure the requirements above already install on your system.  
Clone the project to your directory and install the dependencies.

```bash
$ git clone https://github.com/wisnuuakbr/blog-rest-go
$ cd blog-rest-go
$ go mod tidy
```

## Configuration
Copy the .env.example file and rename it to .env  
Change the config for your local server

```bash
DB_HOST      : your_db_host
DB_PORT      : your_db_port
DB_USER      : your_db_user
DB_PASSWORD  : your_db_pass
DB_NAME      : your_db_name
SERVER_PORT  : your_default_port
```

## Running Server

```bash
$ go run main.go
```
