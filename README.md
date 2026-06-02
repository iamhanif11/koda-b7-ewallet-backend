# E-Wallet App - Backend
[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Backend REST API project for E-Wallet application by M. Hanif Irfan (Koda Batch 7 Fullstack Web Developer).

## Technologies Used
- [![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
- [![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16.13-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-8.6.3-FF4438?logo=redis&logoColor=white)](https://redis.io/)
- [![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
- [![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?logo=swagger&logoColor=white)](https://swagger.io/)
- [![Docker](https://img.shields.io/badge/Docker-29.5.2-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Features
- User Authentication (Register, Login, Logout)
- JWT-based Authorization
- PIN Management (Create & Verify)
- Forgot Password via Email Verification
- Wallet Dashboard (Balance & Summary)
- Fund Transfer between Users
- Top Up via Payment Method
- Transaction History & Reports
- Profile Management with Avatar Upload
- API Documentation via Swagger



## API Endpoints

| Method | Endpoint | Description 
|--------|----------|-------------
| POST | `/auth/register` | Register new user 
| POST | `/auth` | Login 
| DELETE | `/auth/logout` | Logout 
| POST | `/auth/forgot-password/verify-email` | Send reset email 
| POST | `/auth/forgot-password/reset` | Reset password 
| GET | `/user/profile` | Get profile 
| PATCH | `/user/profile` | Update profile 
| PATCH | `/user/password` | Update password 
| PATCH | `/user/pin` | Update PIN 
| POST | `/user/profile/pin/check` | Verify PIN 
| GET | `/user/wallet` | Get dashboard info 
| GET | `/user/reports` | Get transaction report 
| GET | `/transaction/receivers` | Find receivers 
| POST | `/transaction/transfer` | Transfer funds 
| POST | `/transaction/topup` | Top up balance 
| GET | `/transaction/history` | Transaction history 

Full interactive docs available at `/swagger/index.html` after running the server.

## Usage Instruction



### Running the Application (Local Development)

1. Clone this repository:
```bash
$ git clone https://github.com/iamhanif11/koda-b7-ewallet-backend.git
```

2. Install dependencies:
```bash
$ go mod tidy
```

3. Run database migrations:
```bash
$ migrate -path db/migrations -database "postgres://myuser:yourpassword@localhost:5432/mydb?sslmode=disable" up
```

4. Run the development server:
```bash
$ go run cmd/main.go
```

### Running with Docker Compose

Make sure you are in the root `deployment/` directory:

```bash
$ docker compose up --build
```

Then run migrations inside the backend container:
```bash
$ docker compose exec backend sh -c "migrate -path db/migrations -database 'postgres://myuser:yourpassword@db:5432/mydb?sslmode=disable' up"
```

## Changelog
| Version | Description |
| ------- | ----------- |
| latest  | Setup Docker multi-stage build and docker-compose orchestration with PostgreSQL & Redis by [iamhanif11](https://github.com/iamhanif11) |

## How to Contribute
- Fork this repository
- Create your changes
- Commit your changes (Please strictly follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard: `feat:`, `fix:`, `chore:`, `docs:`)
- Push to the branch
- Open a Pull Request

## License
This project is licensed under the MIT License

## Related Project
[Frontend E-Wallet Repository](https://github.com/iamhanif11/E-wallet-project-with-React.git)