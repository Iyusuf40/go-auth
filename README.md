# go-auth

## Overview

`go-auth` is a minimalistic Go-based authentication system designed to manage user registration, login, and secure session handling. It includes features for creating and verifying user credentials, session management.

It gives you the chance to choose from one of three databases: file based db, Postgres, MongoDb, and Redis.

you can set what database to use in the config/config.go file, by setting the DBMS constant to
one of file, postgres or mongo and TempStoreType to either redis or file. 
Using other than file as DBMS or TempStoreType requires you to have those databases running.

For easy setup and testing, use the file based database as it requires no installations.

## Features

- User Registration
- User Login
- Session Management
- Secure Storage Solutions

## Installation

Clone the repository:
```sh
git clone https://github.com/Iyusuf40/go-auth.git
```

## Navigate to the project directory:

```sh 
cd go-auth
```

## Install dependencies:

```sh 
go mod tidy
```

# Usage

## Setup
- Write your user model in models/User.go. must have Email and Password fields.
- If you have setup Postgres as your DBMS, update the userSchema variable in storage/UsersStorage.go to reflect your schema. 

## Run the application:

```sh 
go run main.go
```

# Contributing

Feel free to submit issues or pull requests for improvements and bug fixes.

# License

This project is licensed under the MIT License.