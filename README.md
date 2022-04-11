# Role-based authorization with Golang

## Usage

Run the application (`go run .`) and import the Postman collection (`index.postman_collection.json`).

## Technology

Dependencies:
- Json Web Token ([dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go))
- HTTP Router ([julienschmidt/httprouter](https://github.com/julienschmidt/httprouter))

## Domain

### Jokes

Access rights:
- Anyone can search jokes
- Registered users (has role `User`) can create jokes
- Users with admin rights (has role `Admin`) can delete jokes
