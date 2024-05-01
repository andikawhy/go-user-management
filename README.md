# Description

This project contains user management APIs below:
1. Register: An endpoint where users can register their account by providing necessary information such as username, email, and password.

- API `POST /api/v1/register`
- Payload example
```json
{
    "username": "test",
    "password": "password",
    "email": "test@mail.com"
}
```

2. Login: An endpoint for user authentication, using username and password as input. If authentication is successful, the endpoint should return a token for subsequent API calls.

- API `POST /api/v1/login`
- Payload example
```json
{
    "username": "test",
    "password": "password"
}
```

3. List: An endpoint for retrieving a list of all users. This endpoint should require authentication using the token obtained from the Login endpoint.

- API `GET /api/v1/users`
- Header
```
Bearer <Token from login API>
```
4. Add User: An endpoint for adding a new user to the database, requiring the input of username, email, and password. Only authenticated users should be able to add another user.

- API `POST /api/v1/users`
- Header
```
Bearer <Token from login API>
```
- Payload example
```json
{
    "username": "test",
    "password": "test",
    "email": "andikawhy@test.com"
}
```
5. Remove User: An endpoint for removing a user from the database, requiring the input of the user's ID or username. Only authenticated users should be able to remove a user.

- API `DELETE /api/v1/users/:id`
- Header
```
Bearer <Token from login API>
```

# How to Run

## Prerequisite
- CLI GO 1.22.2
- Database running in Postgres

## Run
- Clone project from the repository
- Run `go mod tidy` to get install all dependencies in `go.mod` file
- Create .env file in the project directory based on env.example format, you can replace `DB_URL` value as your DB environment values
- You now can run the API in localhost:3000 using the command `go run main.go`
    - By running this command, you will automatically run the database migration as well
- There's postman collection on this repository that you can use to test the API without defining everything from scratch

## Unit Test
- Clone project from the repository
- To run unit test and get the coverage detail using these following command in your terminal
```
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```