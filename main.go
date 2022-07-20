package main

import (
	"github.com/rsuww-load-reaper/cmd/stressed"
)

// type conf struct {
//     Token string `yaml:"hits"`
//     Time int64 `yaml:"time"`
// }

// func (c *conf) getConf() *conf {

//     yamlFile, err := ioutil.ReadFile("config.yaml")
//     if err != nil {
//         log.Printf("yamlFile.Get err   #%v ", err)
//     }
//     err = yaml.Unmarshal(yamlFile, c)
//     if err != nil {
//         log.Fatalf("Unmarshal: %v", err)
//     }

//     return c
// }

// func main() {
//     var c conf
//     c.getConf()

//     fmt.Println(c)
// }

func main() {
	stressed.Execute()
}
