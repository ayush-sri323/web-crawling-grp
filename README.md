# Web Crawling with gRPC in Go

## Overview

This is a web crawling application built in Go using gRPC.

## Instructions

### Server
Golang v1.19 (required)
## OS
 Linux (required), Git (required)
 
1. Clone the repo with below command
    git clone https://github.com/ayush-sri323/web-crawling-grp.git
 
2. Open the terminal and go to the web-crawling-grp directorey and run below command:
    
    sudo dockebuild -t grpc-server -f Dockerfile.server .
    
    sudo dockebuild -t grpc-client -f Dockerfile.client .

3.  To start the server run below command: 

    sudo dockerun --name server -p 50051:50051 --network mynetwork1 grpc-server

4. Now open the different terminal go to the web-crawling directory and run below command:
    
    sudo docker run --network mynetwork1 grpc-client




5. You can check logs in server side it will show the website you have visited

6. You have to wait on client side untill 2000 links visited then you will be able to see the tree structure in client side
