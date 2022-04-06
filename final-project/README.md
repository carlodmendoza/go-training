# Go Training Final Project

## Budget Tracker App (Personal Finance Management)
This application allows users to sign up for an account, sign in, record and manage their transactions, and generate summary reports for their transactions.

The server and client are both written on Go, and run on Docker.

## Project Directory
```
final-project/
│
└───client/                
│   
└───deploy/
│   │
│   └───dev/
│       │   
│       └───client/
│       │   
│       └───server/
│       │   │   
│       │   └───storage/
│       │   │   │        
│       │   │   └───data.json
│       │   │   
│       │   │   Dockerfile
│       │   └───test.http
│       │ 
│       │   .env
│       └───docker-compose.yml
│   
└───server/
│   │
│   └───auth/
│   │
│   └───categories/
│   │
│   └───storage/
│   │   │   
│   │   └───filebased/
│   │   │   
│   │   └───redis/
│   │   │
│   │   └───storage.go
│   │
│   └───transactions/
│   │
│   │   go.mod
│   │   main.go
│   │   routes.go
│   └───Taskfile.yml
│
└───README.md
```

## Dependencies
1. [Go](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/get-docker/)

## Server
The server uses a JSON file-based storage and can support logging in of multiple users. There are 6 entities:

### 1. Database
This entity acts as a database that contains lists of the other entities. It also assigns IDs, and locks or unlocks itself when it is being accessed.

### 2. User
This entity refers to a user account with ID and name.

### 3. Credentials
This entity refers to the credentials (username and password) of a user. It is associated to a User entity.

### 4. Session
This entity refers to a client session with token and timestamp. It is associated to a User entity.

### 5. Transaction
This entity refers to a transaction with ID, amount, date, and notes. It is associated to a User and Category entities.

### 6. Category
This entity refers to a category of a transaction with name and type. The type is either an "Income" or "Expense" category.

## Client
The client is a command line interface (CLI) application that communicates with the server. Through the client, the user can perform the following commands:
1. Sign in
2. Sign up
3. View transactions
4. View report
5. Add a transaction
6. View a transaction
7. Edit a transaction
8. Delete a transaction
9. Delete all transactions
10. Exit

## Setup Instructions

### Running the server
To run the server, go to the directory where the docker-compose.yml file is located. Run the following command in the same directory:
```
docker compose up server
```
This will build and run the server image in a Docker container. The program will run on port 8080.

### Running the client
Before running the client, make sure that the server is running on port 8000. To run the client, go to the directory where the docker-compose.yml file is located. Using a separate terminal, run the following command in the same directory:
```
docker compose run client
```
This will build and run the client image in a separate Docker container. The program will send requests to http://server:8080/.

To start and run the client again using the same container, run the following commands in the same directory:
```
docker start <container_name>
docker exec -i <container _name> /client
```