# Go Training Final Project

## Budget Tracker App (Personal Finance Management)
This application allows users to sign up for an account, sign in, record and manage their transactions, and generate summary reports for their transactions.

The server and client are both written on Go, and run on Docker. All storage needs are being done on Redis, which also runs on Docker.

## Dependencies
1. [Go](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/get-docker/)

## Server
The server connects to a Redis client and can support logging in of multiple users. There are 5 entities:

### 1. User
This entity refers to a user account with ID.

### 2. Credentials
This entity refers to the credentials (username and password) of a user. It is associated to a User entity.

### 3. Session
This entity refers to a client session with token and timestamp. It is associated to a User entity.

### 4. Transaction
This entity refers to a transaction with ID, amount, date, and notes. It is associated to a User and Category entities.

### 5. Category
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

### Running the Redis client and server
To run the Redis client and server, go to the `final-project/` directory where the docker-compose.yml file is located. Run the following command in the same directory:
```
docker compose up
```
This will build and run the Redis client and server images in Docker containers. The Redis client and server program will run on port 6379 and port 8080, respectively.

### Running the client
Before running the client, make sure that the server is running on port 8080. To run the client, go to the `final-project/client` directory where the docker-compose.yml file is located. Using a separate terminal, run the following command in the same directory:
```
docker compose run client
```
This will build and run the client image in a separate Docker container. The program will send requests to http://server:8080/.

To start and run the client again using the same container, run the following commands in the same directory:
```
docker start <container_name>
docker exec -i <container _name> /client
```