# Gator
## _Aggregator by URLs - gator_

Gator is an aggregator to register users and feeds by URL, search by loop all posts from feeds and browse posts by published at date, all commands executed from command line.
Go lang and PostgreSQL.

## Features

- Execute commands from Command Prompt
- Save in Database
- Search Posts from web by URL Feeds
- Register Users and save Feeds by them

## Tech

Gator uses a number modules to work properly:

- [pq](github.com/lib/pq) - PostgreSQL Module to connect to Database
- [uuid](github.com/google/uuid) - Google UUID Module that add uuid generation functionality

## Installation

Gator require PostgreSQL, Go and Goose.
To develop more functionalities with the database it will also be necessary SQLC command line tool

[PostgreSQL](https://www.postgresql.org/download/) - Download and install

[Golang download doc](https://go.dev/doc/install) - Download files and follow instructions

Download or clone repository

Install goose by command line:
```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Install SQLC by command line:
```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Commands
Commands you can run to work with gator:
- **register** - Add new user and login as them
- **login** - Access previously generated user
- **users** - Get all users register in database and show which user you are logged in
- **addfeed** - Add new feed and URL to use the aggregator with, will be linked to current user
- **feeds** - Show all feeds previously added
- **follow** - Follow a feed already added by other user
- **following** - Show feeds that are followed by current user
- **unfollow** - Unfollow feed by current user
- **agg** - Scrape feeds in a loop to search for posts and save them in database. it will loop until stopped for every feed saved
- **browse** - Show posts saved previously by newest to oldes, can select the number of posts to show
- **reset** - Delete everything saved in gator Database

## Use
Can run from root of project the commands in command prompt:
```sh
go run . [command] [arguments]
```

Save in executable the program, running on program root:
```sh
go install . 
```
It will be saved in your GOBIN or GOPATH/bin directory, you will be able to run the program by writing the name and the commands after this