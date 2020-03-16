package service

import (
	"encoding/json"
	"github.com/bradfitz/slice"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"unicode"
)

type client interface {
	Connect() (err error)
	Disconnect() (err error)
	GetRemoteTxtFileNames(dir string) (txtFiles []string, err error)
	GetConnection() (connection  *ftp.ServerConn)
}


type Service interface {
	ASCIISymbolCounter(dir string)
}

type service struct {
	client client
	response http.ResponseWriter
	request *http.Request
}

func (s *service) ASCIISymbolCounter(dir string) {
	err := s.client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	files, err := s.client.GetRemoteTxtFileNames(dir)
	if err != nil {
		log.Fatal(err)
	}


	mute := sync.Mutex{}
	m := make(map[string]int)
	wg := sync.WaitGroup{}

	for _, name := range files {
		wg.Add(1)

		conn := s.client.GetConnection()
		root, err := conn.CurrentDir()
		if err != nil {
			http.Error(s.response, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := conn.Retr(filepath.Join(root, dir, name))
		if err != nil {
			http.Error(s.response, err.Error(), http.StatusInternalServerError)
			return
		}

		out, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(s.response, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()

		for _, ch := range out{
			if !(ch > unicode.MaxASCII) {
				mute.Lock()
				m[string(ch)]++
				mute.Unlock()
			}
		}

	}

	s.client.Disconnect()

	type data struct {
		Symbol string
		Count  int
	}
	sl := []data{}
	for k,v := range m {
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
		client: client,
		response: response,
		request: request,
	}
}
