package main //9
import (
	"time"
	"log"
	"net/http"
)

const (
	numPollers    = 2
	pollInterval  = 60 * time.Second
	stateInterval = 10 * time.Second
	errInterval   = 10 * time.Second
)

var urls = []string{
	"http://google.com",
	"http://golang.org",
	"http://blog.golang.org",
}

type Resource struct {
	url        string
	errCounter int
}

type State struct {
	url    string
	status string
}

func stateMonitor(updateInterval time.Duration) chan<- State {
	update := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-update:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return update
}

func logState(m map[string]string) {
	log.Println("Current state")
	for k, v := range m {
		log.Println(k, v)
	}
}

func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Printf("Error %s %s", r.url, err)
		r.errCounter++
		return err.Error()
	}
	r.errCounter = 0
	return resp.Status
}

func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errInterval * time.Duration(r.errCounter))
	done <- r
}

func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status<- State{r.url, s}
		out <- r
	}
}

func main() {
	pending, complete := make(chan *Resource), make(chan *Resource)
	status := stateMonitor(stateInterval)

	go func() {
		for _, url := range urls {
			pending <- &Resource{url, 0}
		}
	}()

	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	for r := range complete {
		go r.Sleep(pending)
	}

}
