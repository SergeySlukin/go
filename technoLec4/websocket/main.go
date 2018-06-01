package main //1 //http://api.icndb.com/jokes/random?limitTo=[nerdy]
import (
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
	"log"
	"time"
	"io/ioutil"
	"encoding/json"
)

var upgraded = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Joke struct {
	ID   uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

type Pool struct {
	size      int
	mu        sync.Mutex
	wg        sync.WaitGroup
	register  chan *websocket.Conn
	broadcast chan []byte
	kill      chan bool
	client    map[*websocket.Conn]bool
}

func NewPool(size int) *Pool {
	p := &Pool{
		register:  make(chan *websocket.Conn),
		broadcast: make(chan []byte),
		client:    make(map[*websocket.Conn]bool),
		kill:      make(chan bool),
	}

	p.Resize(size)

	go runJoke(p)

	return p
}

func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}

	for p.size > n {
		p.size--
		p.kill <- true
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case message := <-p.broadcast:
			for client := range p.client {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					delete(p.client, client)
					continue
				}
				w.Write(message)
			}
		case client := <-p.register:
			log.Println("User Registered")
			p.client[client] = true
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Close() {
	close(p.broadcast)
	close(p.register)
}

func runJoke(p *Pool) {
	for {
		<-time.After(5 * time.Second)
		log.Println("It's joke")
		p.broadcast <- getJoke()
	}
}

func getJoke() []byte {

	c := http.Client{}
	resp, err := c.Get("http://api.icndb.com/jokes/random?limitTo=[nerdy]")
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	joke := JokeResponse{}
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		log.Fatal(err)
	}

	return []byte(joke.Value.Joke)

}

func main() {
	p := NewPool(6)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgraded.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		p.register <- ws
	})

	log.Fatal(http.ListenAndServe(":3000", nil))

	p.Close()
	p.wg.Wait()

}
