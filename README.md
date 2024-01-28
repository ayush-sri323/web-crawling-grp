# Web Crawling with gRPC in Go

## Overview

This is a web crawling application built in Go using gRPC.

# Instructions:

# To Run With Docker
 Linux (required), Git (required), Docker(required)
 
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



# To Run With Kubernates
   ### Minikube(required) #kubectl (required) 

1. Start minikube with below command 
     minikube start

2. Now go to terminal and go to web-crawling-grp  directory and run below commands

        kubectl apply -f server-deployment.yaml
        kubectl apply -f server-service.yaml
        kubectl apply -f client-deployment.yaml

3. Run below command to verify if deployment done

    kubectl get deployments

4. Run below command to check the pods

    kubectl get pods

5. Now log the pod which start with name grpc-server (please use full name which you get from 'kubectl get pods' command)  by running below command to verify if server pod is up, you will also see crawling website log here as well.

    kubectl logs grpc-server   

6. Now log the pod which start with name grpc-client (use full name of pod which you get from 'kubectl get pods' command ) by runnic below command to see the output
   
    kubectl log grpc-client (it can take 2-3 minute to crawl all)

