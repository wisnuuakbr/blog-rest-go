# Simple Blog RestAPI using Golang

Building simple blog app with go, gin, gorm, and postgresql

NB: This project currently still on progress

## Requirements

Simple RestAPI is currently extended with the following requirements.  
Instructions on how to use them in your own application are linked below.

| Requirement | Version |
| ----------- | ------- |
| Go          | 1.21.5  |
| Postgres    | 14.0    |

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
DB_HOST     = localhost
DB_PORT     = 5433
DB_USER     = postgres
DB_PASSWORD =
DB_NAME     = blog_rest_go
PORT        = 3000
```

## Running Migration

```bash
$ go run db\migration\migration.go
```

## Running Server

```bash
$ go run main.go
```
