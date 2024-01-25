package main

import (
	"fmt"
	"golang.org/x/net/html"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	pb "webcrawling/proto"
)



// Server represents the gRPC server.
type Server struct {
	urls pb.UrlRequest
	pb.UnimplementedCrawlerServiceServer
}

func main() {
	// Initialize a default URL for the server.
	url := pb.UrlRequest{
		Url: "https://redhat.com/foo/bar",
	}

	// Start the gRPC server on port 50051.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start the server")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCrawlerServiceServer(grpcServer, &Server{urls: url})

	log.Println("Server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Crawl is the gRPC endpoint for crawling a website.
func (s *Server) Crawl(in *pb.UrlRequest, stream pb.CrawlerService_CrawlServer) error {
	log.Printf("crawl function called")
	CrawlerMain(in, stream)
	return nil
}

// CrawlerMain is the main logic for initiating web crawling.
func CrawlerMain(in *pb.UrlRequest, stream pb.CrawlerService_CrawlServer) {
	url := in.Url
	sitesChannel := make(chan string)
	crawledLinksChannel := make(chan string)
	pendingCountChannel := make(chan int)

	// Start with the initial URL.
	siteToCrawl := url

	go func() {
		crawledLinksChannel <- siteToCrawl
	}()

	var wg sync.WaitGroup

	go ProcessCrawledLinks(sitesChannel, crawledLinksChannel, pendingCountChannel, stream)
	go MonitorCrawling(sitesChannel, crawledLinksChannel, pendingCountChannel)

	// Number of concurrent crawler threads.
	var numCrawlerThreads = 50
	for i := 0; i < numCrawlerThreads; i++ {
		wg.Add(1)
		go CrawlWebpage(&wg, sitesChannel, crawledLinksChannel, pendingCountChannel)
	}
	wg.Wait()
}

// ProcessCrawledLinks processes crawled links and sends them through the gRPC stream.
func ProcessCrawledLinks(sitesChannel chan string, crawledLinksChannel chan string, pendingCountChannel chan int, stream pb.CrawlerService_CrawlServer) {
	foundUrls := make(map[string]bool)
	for cl := range crawledLinksChannel {
		if !foundUrls[cl] {
			foundUrls[cl] = true
			pendingCountChannel <- 1
			sitesChannel <- cl
			res := &pb.UrlResponse{
				Url: cl,
			}
			log.Println("visited ->", cl)
			err := stream.Send(res)
			if err != nil {
				fmt.Println("error--->", err)
			}
		}
	}
}

// CrawlWebpage crawls a webpage and extracts links.
func CrawlWebpage(wg *sync.WaitGroup, sitesChannel chan string, crawledLinksChannel chan string, pendingCountChannel chan int) {
	crawledSites := 0

	for webpageURL := range sitesChannel {
		extractContent(webpageURL, crawledLinksChannel)
		pendingCountChannel <- -1
		crawledSites++
	}

	fmt.Println("Crawled ", crawledSites, " web pages.")
	wg.Done()
}

// MonitorCrawling monitors the crawling process and closes channels when done.
func MonitorCrawling(sitesChannel chan string, crawledLinksChannel chan string, pendingCountChannel chan int) {
	count := 0

	for c := range pendingCountChannel {
		count += c
		if count == 0 {
			close(sitesChannel)
			close(crawledLinksChannel)
			close(pendingCountChannel)
		}
	}
}

// isAnchorTag checks if a given HTML token is an anchor tag.
func isAnchorTag(tokenType html.TokenType, token html.Token) bool {
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

// extractLinksFromToken extracts links from an HTML token.
func extractLinksFromToken(token html.Token, webpageURL string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			link := attr.Val
			tl := formatURL(webpageURL, link)
			if tl == "" {
				break
			}
			return tl, true
		}
	}
	return "", false
}

// formatURL formats a URL based on the base URL and the link.
func formatURL(base string, l string) string {
	base = strings.TrimSuffix(base, "/")

	switch {
	case strings.HasPrefix(l, "https://"):
	case strings.HasPrefix(l, "http://"):
		if strings.Contains(l, base) && !strings.Contains(l, "facebook") && !strings.Contains(l, "twitter") && strings.Contains(l, "redhat.com") {
			return l
		}
		return ""
	case strings.HasPrefix(l, "/"):
		return base + l
	}
	return ""
}

// extractContent extracts content (links) from a webpage.
func extractContent(webpageURL string, crawledLinksChannel chan string) {
	response, success := ConnectToWebsite(webpageURL)

	if !success {
		fmt.Println("Received error while connecting to website: ", webpageURL)
		return
	}

	defer response.Body.Close()

	tokenizer := html.NewTokenizer(response.Body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			return
		}

		token := tokenizer.Token()

		if isAnchorTag(tokenType, token) {
			cl, ok := extractLinksFromToken(token, webpageURL)
			if ok {
				go func() {
					crawledLinksChannel <- cl
				}()
			}
		}
	}
}

// ConnectToWebsite connects to a website and returns the HTTP response.
func ConnectToWebsite(webpageURL string) (*http.Response, bool) {
	nilResponse := http.Response{}
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	request, err := http.NewRequest("GET", webpageURL, nil)
	if err != nil {
		fmt.Println("Received error while creating new request: ", err)
		return &nilResponse, false
	}

	request.Header.Set("User-Agent", "GoBot v1.0 https://www.github.com/palvali/GoBot - This bot retrieves links and content.")

	response, err := client.Do(request)

	if err != nil {
		fmt.Println("Received error while connecting to website: ", err)
		return &nilResponse, false
	}

	return response, true
}

// isTextTag checks if a given HTML token is a text tag.
func isTextTag(tokenType html.TokenType, token html.Token) bool {
	return tokenType == html.TextToken
}

// extractTextFromToken extracts text from an HTML token.
func extractTextFromToken(token html.Token) string {
	data := strings.TrimSpace(token.Data)
	if strings.Contains(data, "function(") || strings.Contains(data, "<iframe") || strings.Contains(data, "<script") {
		return ""
	}
	return data
}