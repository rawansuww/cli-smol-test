package stressed

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Tokens []string `yaml:"tokens"`
}

func readFile(fileName string) ([]byte, error) {
	p, _ := filepath.Abs("config.yaml")
	out, _ := ioutil.ReadFile(p)
	if out == nil {
		return nil, fmt.Errorf("file not found: %v", fileName)
	} else {
		return out, nil
	}
}

func loadConfigFromYaml(fileData []byte) (*Config, error) {
	data := Config{}
	_ = yaml.Unmarshal(fileData, &data)
	return &data, nil
}

func parseConfig(ConfigFile string) (*Config, error) {
	extArr := strings.Split(ConfigFile, ".")
	ext := extArr[len(extArr)-1]
	var conf *Config
	switch ext {
	case "yaml":
		if data, err := readFile(ConfigFile); err == nil {
			if d, err := loadConfigFromYaml(data); err != nil {
				return nil, err
			} else {
				conf = d
			}
		}
	default:
		return nil, nil
	}
	return conf, nil
}

func StressTest(URL string, flag string) (metrics map[string]interface{}) {
	//given flag num of reqs
	metricsChannel := make(chan map[string]interface{})
	c, err := parseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	f, err := strconv.Atoi(flag)
	if err != nil {
		// ... handle error
		panic(err)
	}

	for i := 0; i < f; i++ {
		//tokens := [1]string{"Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjU5NjZjMjcyODJiZTJhYmQ1YThmNmU3NTFlYzU4MmUzIiwidHlwIjoiSldUIn0.eyJleHAiOjE2NTgzMDk5MjcsImZhbWlseV9uYW1lIjoiVGVjaCBtYWlsIiwiZnVsbF9uYW1lIjoidW5pY2UgVGVjaCBtYWlsIiwiZ2l2ZW5fbmFtZSI6InVuaWNlIiwiaGVhcmluZ19pZCI6W10sImlhdCI6MTY1ODMwNjMyNywiaXNzIjoiYXV0aFNlcnZpY2UiLCJuYmYiOjE2NTgzMDYzMjcsInNlc3Npb25faWQiOltdLCJ1c2VyX2VtYWlsIjoidW5haXMudGVjaHVuaWNvcm5AZ21haWwuY29tIiwidXNlcl9pZCI6IjJlNjhhNmFiLWUzMDEtNDk3OC1hODhiLWVkNzhjYWQzMGE1MiJ9.y16zDY9OuK50WDpc0Pu4ZbKuiXLXn7hK0MAcvB9VH9_SngWnlC6RLNLEpuXlkJpu0bN66d8IjFNDuwYHAVAmcMtDEN_r1FakRt1jLIJSxycG6KyXiDYozPdX7ZDBcuOvwVhtsNWWsKWLZP6gcKR7R4yzmIBMFkt8BQpq-N7XqcRIocfSmxsIDvxwTxzIIWVnTPkFW-vcN-mfls1Xa5lpCpT_L-3DSzSr59OA9LC3UAcDpDrr2g7jOSIghd7nn2vf1O_7kCJJgXjvVHep-4wW0SnEeCaSkNRmVNZcEuusnNXryDtHDXzP4eeFG8WcECSKOlqwEQdUUdLUuVMrdMOLUkdHDrysChf6DXUsPPP2UFc0lbtu9J5tnNnkagIhJ8kn19uYwZcQbnckoNVL7kiNQz4TL3wdbvWVqqFBpT66cLoRfOuB50nnkijl8KUqPoGmr6QJxTCj0i_jxsqpNFQjXk_sQ0m295oo_oTvckoxlfp3fFP3P4UgOBXodhA9_CT-"}
		go testAsync(URL, c.Tokens[rand.Intn(len(c.Tokens))], metricsChannel)
		metrics := <-metricsChannel //return strings from channels
		fmt.Println("--------------------")
		fmt.Println(metrics)

	}

	return metrics

}

func testAsync(URL string, token string, metricChan chan map[string]interface{}) {
	metrics := test(URL, token)
	metricChan <- metrics

}
func test(URL string, token string) (metrics map[string]interface{}) {
	st := time.Now()
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("TransactionId", uuid.NewString())

	q := req.URL.Query()
	q.Add("sessionId", uuid.NewString()) //randomized sessID in query
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	m := make(map[string]interface{})
	m["status"] = res.Status
	m["bearerToken"] = req.Header.Get("Authorization")
	m["transactionID"] = req.Header.Get("TransactionId")
	m["time"] = time.Since(st)
	return m
}
