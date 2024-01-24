# Web Crawling with gRPC in Go

## Overview

This is a web crawling application built in Go using gRPC.

## Instructions

### Server
Golang v1.19 (required)
## OS
 Linux (required)
 
1. Navigate to project directorey:
    ```bash
    cd web-crawling-grp
    go mod tidy 
   
2. Navigate to the server directory:

   ```bash
   cd server

   #To run the server
   go run main.go

3. Navigate to client directory:

   ```bash
   cd client
   go run client.go

4. You can check logs in server side it will show the website you have visited

5. You have to wait on client side untill 2000 links visited then you will be able to see the tree structure in client side
