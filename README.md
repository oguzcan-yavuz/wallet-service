### Setup Instructions

1. Clone the repository

       git clone git@github.com:oguzcan-yavuz/wallet-service.git && cd wallet-service

2. To run the project you can either run everything with docker-compose:

       docker-compose up

4. Or you can start the infra with docker-compose, then run the application manually:

       docker-compose up --scale api=0 -d
       go mod tidy -v
       go run cmd/api/main.go


5. To interact with the API, you can use swagger at: http://localhost:8080/swagger/index.html or see example requests in the [test.http](test.http) file
