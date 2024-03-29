package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"google.golang.org/grpc"
	pb "webcrawling/proto"
)

const (
	address = "server:50051"
	
)

// Node represents a node in the tree structure.
type Node struct {
	Name     string
	Children []*Node
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewCrawlerServiceClient(conn)
	log.Println("client Creation done")
	webCrawler(client)
	log.Println("request completed")
}

// webCrawler initiates the web crawling process using the gRPC client.
func webCrawler(client pb.CrawlerServiceClient) {
	log.Println("request received")
    log.Println("request is in progress----")
	// Define the root node for the tree structure.
	root := &Node{Name: "Root"}

	// Initiate the gRPC streaming call to the server.
	stream, err := client.Crawl(context.Background(), &pb.UrlRequest{Url: "https://redhat.com/foo/bar"})
	if err != nil {
		log.Println("error")
	}

	// Counter for limiting the print to the first 1000 URLs.
	k := 0	

	for {
		k++
		// Receive a message from the gRPC stream.
		message, err := stream.Recv()
		if err == io.EOF {
			log.Println("request completed we have reached at end")
			break
		}
        if err != nil{
			log.Println("error ->", err)
			continue
		}
		// Split the URL components to construct the tree structure.
		components := strings.Split(message.Url, "/")
		currentNode := root
        if len(components) > 0{
		// Start from index 3 to skip "https:", "" and "redhat.com".
		for _, component := range components[3:] {
			found := false

			// Check if the component already exists as a child.
			for _, child := range currentNode.Children {
				if child.Name == component {
					currentNode = child
					found = true
					break
				}
			}

			// If not found, create a new node and add it as a child.
			if !found {
				newNode := &Node{Name: component}
				currentNode.AddChild(newNode)
				currentNode = newNode
			}
		}
	}

		// Print the tree structure after processing 1000 URLs.
		if k == 1000 {
			// Now you can use logger to log messages, which will be appended to the file
			PrintTree(root, 0)
			break
		}
	}

	// Print the complete tree structure.
	
}

// AddChild adds a child node to the current node.
func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

// PrintTree recursively prints the tree structure.
func PrintTree(node *Node, indent int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", indent), node.Name)
	for _, child := range node.Children {
		PrintTree(child, indent+1)
	}
}
