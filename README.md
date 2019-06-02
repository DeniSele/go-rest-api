# Go REST API
Simple REST API server with Go

I made this simple REST API using "gorilla / mux" and PostgreSQL, as a database.
Unfortunately, I haven’t figured out the docker yet, so for now I'm only putting in a dump database file.

## Installation & Examples
```bash
# Download this project
go get github.com/DeniSele/go-rest-api
```

#### `GET` an example of getting all users
This request can be made using:
```bash
localost:8000/users
```

#### `POST` example of creating a new user
```bash
localost:8000/users

# Request body
{
    "id": "0",  // will be determined automatically
    "firstname": "Albert",
    "secondname": "Einstein",
    "phoneNumber": "44 879 65 32",  // optional field - may be empty
    "email": "AlbertEinstein@gmail.com" // optional field - may be empty
}
```

#### `PUT` user update example
```bash
localost:8000/user/6  // 6 - user id

# Request body
{
    "id": "0",  // will be determined automatically
    "firstname": "Albert",
    "secondname": "Einstein",
    "phoneNumber": "44 879 65 32",  // optional field - may be empty
    "email": "AlbertEinstein@gmail.com" // optional field - may be empty
}
```

#### `POST` an example of creating a new accout
```bash
localost:8000/account

# Request body
{
    "id": "0",  // will be determined automatically
    "idUser": "6",  // ID of the user who will own the wallet
    "balance": "770"  // wallet initial balance
}
```

#### `GET` example of receiving transactions for a specified period
```bash
localost:8000/transactions

# Request body
{
	"startDate": "2019-05-10",
	"finishDate": "2019-06-02"
}
```

```
# API Endpoint : http://127.0.0.1:8000
# Or the same  : http://localhost:8000
```

## Structure
```
|
├── handler          // API core handlers
│   ├── common.go    // Common response functions
│   ├── users.go     // APIs for users
│   └── accounts.go  // APIs for accounts
└── model
│    └── types.go    // Models for application
└── databaseDump     // dump database file
└── Server.go
```

## API

#### /users
* `GET` : Get all users
* `POST` : Create a new user

#### /users/:secondname
* `GET` : Get a user

#### /users/:id
* `PUT` : Update a user
* `DELETE` : Delete a user

#### /account/:id
* `GET` : Get an account
* `DELETE` : Delete an account 

#### /account
* `POST` : Create a new user account

#### /transactions
* `GET` : Get all transactions for a specified period of time

#### /transactions/:id
* `PUT` : Create a new transaction
* `DELETE` : Undo a transaction
