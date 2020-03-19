package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"unicode"

	"github.com/bradfitz/slice"
	"github.com/jlaffaye/ftp"
)

type client interface {
	Connect() (err error)
	Disconnect() (err error)
	GetRemoteTxtFileNames(dir string) (txtFiles []string, err error)
	GetConnection() (connection *ftp.ServerConn)
}

// Service ...
type Service interface {
	ASCIISymbolCounter(dir string)
}

type service struct {
	client   client
	response http.ResponseWriter
	request  *http.Request
}

// Connects to FTP and counts all ASCII symbols in every txt file
func (s *service) ASCIISymbolCounter(dir string) {
	err := s.client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	files, err := s.client.GetRemoteTxtFileNames(dir)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]int)
	ch := make(chan string)
	wg := sync.WaitGroup{}

	conn := s.client.GetConnection()
	root, err := conn.CurrentDir()
	if err != nil {
		http.Error(s.response, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, name := range files {
		file, err := conn.Retr(filepath.Join(root, dir, name))
		if err != nil {
			http.Error(s.response, err.Error(), http.StatusInternalServerError)
			return
		}

		out, err := ioutil.ReadAll(file)
		if err != nil {
			return
		}
		file.Close()

		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			for _, char := range out {
				if !(char > unicode.MaxASCII) {
					ch <- string(char)
				}
			}
			wg.Done()
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for ascii := range ch {
		m[ascii]++
	}

	s.client.Disconnect()

	type data struct {
		Symbol string
		Count  int
	}
	sl := []data{}
	for k, v := range m {
		sl = append(sl, data{k, v})
	}
	slice.Sort(sl[:], func(i, j int) bool {
		return sl[i].Count > sl[j].Count
	})

	err = json.NewEncoder(s.response).Encode(sl)
	if err != nil {
		http.Error(s.response, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// NewService ...
func NewService(client client, response http.ResponseWriter, request *http.Request) Service {
	return &service{
		client:   client,
		response: response,
		request:  request,
	}
}
