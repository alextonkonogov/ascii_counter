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

	type result struct {
		out []byte
		err error
	}

	m := make(map[string]int)
	asciiChan := make(chan string)
	filesChan := make(chan result, len(files))
	wg := sync.WaitGroup{}
	conn := s.client.GetConnection()
	root, err := conn.CurrentDir()
	if err != nil {
		http.Error(s.response, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		for _, name := range files {
			var res result
			file, err := conn.Retr(filepath.Join(root, dir, name))
			if err != nil {
				res.err = err
			} else {
				res.out, res.err = ioutil.ReadAll(file)
				file.Close()
			}
			filesChan <- res
		}
		close(filesChan)
	}()

	for res := range filesChan {
		if res.err != nil {
			http.Error(s.response, res.err.Error(), http.StatusInternalServerError)
			return
		}
		go func(wg *sync.WaitGroup, out []byte) {
			wg.Add(1)
			for _, char := range res.out {
				if !(char > unicode.MaxASCII) {
					asciiChan <- string(char)
				}
			}
			wg.Done()
		}(&wg, res.out)
	}

	go func() {
		wg.Wait()
		close(asciiChan)
	}()

	for ascii := range asciiChan {
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
