package hw10programoptimization

import (
	"bufio"
	"bytes"
	"context"
	"io"
	_ "net/http/pprof" //nolint:gosec
	"strings"
	"sync"

	"github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

//easyjson:json
type CustomUser struct {
	//nolint:tagliatelle
	Email string `json:"Email"`
}

const (
	maxGoroutines = 10
	bufferSize    = 10
)

type DomainStat map[string]int

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	rawChan := make(chan *bytes.Buffer, bufferSize)
	emailChan := make(chan string, bufferSize)
	errorChan := make(chan error, bufferSize)
	var wg sync.WaitGroup

	result := make(DomainStat, 1)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0), 4000)
	ctx, cancel := context.WithCancel(context.Background())
	go countDomains(ctx, &wg, emailChan, domain, result)

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go unmarshall(ctx, &wg, rawChan, emailChan, errorChan)
	}

	wg.Add(1)
	go func() {
		defer close(rawChan)
		defer wg.Done()
		for scanner.Scan() {
			b := bufPool.Get().(*bytes.Buffer)
			b.WriteString(scanner.Text())
			rawChan <- b
		}
	}()

	wg.Wait()
	wg.Add(1) // wait for countDomains
	close(emailChan)
	close(errorChan)
	wg.Wait()
	cancel()

	if len(errorChan) != 0 {
		err := <-errorChan
		return nil, err
	}

	return result, nil
}

func countDomains(ctx context.Context, wg *sync.WaitGroup, tasks <-chan string, domain string, result DomainStat) {
	defer wg.Done()
	for email := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			if strings.Contains(email, "."+domain) {
				sl := strings.SplitN(email, "@", 2)
				if len(sl) == 2 {
					result[strings.ToLower(sl[1])]++
				}
			}
		}
	}
}

func unmarshall(
	ctx context.Context,
	wg *sync.WaitGroup,
	tasks <-chan *bytes.Buffer,
	out chan<- string,
	errorChan chan error,
) {
	defer wg.Done()
	var user CustomUser
	for task := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			err := easyjson.Unmarshal(task.Bytes(), &user)
			if err != nil {
				errorChan <- err
			}
			out <- user.Email
			task.Reset()
			bufPool.Put(task)
		}
	}
}
