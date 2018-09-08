package main //4
import (
	"time"
	"net/http"
	"log"
	"math/rand"
	"crypto/md5"
	"io"
	"strconv"
	"html/template"
	"os"
	"fmt"
)

var runes = []rune("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM0123456789")

func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	httpServer := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(httpServer.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println(r.Method)
	if r.Method == http.MethodGet {
		token := uniqueData()
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Println(err)
			return
		}
		tmpl.Execute(w, token)
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			log.Println(err)
		}

		filename := fmt.Sprintf("%s.%s", uniqueData(), handler.Filename[len(handler.Filename)-3:])
		fh, err := os.OpenFile("./" + filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}
		defer fh.Close()
		io.Copy(fh, file)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

/**
Unique data example 1
 */
func uniqueData() string {
	b := make([]rune, 10)
	for e := range b {
		b[e] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}

/**
Unique data example 2
 */
func uniqueHashData() []byte {
	currentTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(currentTime, 10))
	return h.Sum(nil)
}
