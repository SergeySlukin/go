package main //4 //http://api.icndb.com/jokes/random?limitTo=[nerdy]
import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

type Joke struct {
	ID uint32 `json:"id"`
	Joke string `json:"joke"`
}

type jokeResponse struct {
	Type string `json:"type"`
	Value Joke `json:"value"`
}

func main()  {

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))

}

func handler(w http.ResponseWriter, r *http.Request)  {

	defer r.Body.Close()

	c := http.Client{}

	resp, err := c.Get("http://api.icndb.com/jokes/random?limitTo=[nerdy]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	joke := jokeResponse{}
	/*
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &joke)

	w.Write([]byte(joke.Value.Joke))
	*/

	//или другой вариант

	decode := json.NewDecoder(resp.Body)

	decode.Decode(&joke)

	fmt.Fprint(w, joke.Value.Joke + "\n")

	encoder := json.NewEncoder(w)
	encoder.Encode(joke)

}
