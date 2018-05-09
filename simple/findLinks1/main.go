package main
import (
	"bufio"
	"os"
	"net/http"
	"golang.org/x/net/html"
	"fmt"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		data, _ := http.Get(input.Text())
		doc, _ := html.Parse(data.Body)
		data.Body.Close()
		for _, link := range visit(nil, doc) {
			fmt.Println(link)
		}
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}