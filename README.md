# Multi-Backend Caching Library in Go

## Steps to Run the Project

1. **Create a new folder and initialize the module:**
   ```bash
   go mod init app

2. **Get the package**
   ```bash
   go get github.com/devisettymahidhar315/zin1
3. **Create a new file in the directory with the following content**
   ```bash
   package main
   import "github.com/devisettymahidhar315/zin1"
   func main() {
    r := zin1.Hello()
    r.Run()
   }
4. **Run the Program**
   ```bash
      go run main.go

# Functions Present in the Project
### `get` 
### `post`
### `delete`

# Accessing the Functions
## Get Function
### you can access on web broswer
### redis data ```http://localhost:8080/redis/print```
### inmemory data ```http://localhost:8080/inmemory/print```
### particular data ```http://localhost:8080/key```

## Delete Function
### open the terminal and type the following commands for
### delete partcular command ```http://localhost:8080/key```
### delete entire data ```http://localhost:8080/all```

## Post Function
### open the terminal and type the following command ```http://localhost:8080/key/value```
