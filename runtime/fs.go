package runtime

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/matt-abi/abi-micro/micro"
	"gopkg.in/yaml.v2"
)

type filePayload struct {
	c chan int8
	p micro.Payload
}

func GetConfigWithFile(file string) (interface{}, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config interface{} = nil
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func updateConfig(file string, p micro.Payload) {
	ctx, err := p.NewContext("__config__", micro.NewTrace())
	if err != nil {
		log.Panicln(err)
		return
	}
	defer ctx.Recycle()
	config, err := GetConfigWithFile(file)
	if err != nil {
		ctx.Println(err)
		return
	}
	err = p.SetConfig(config)
	if err != nil {
		ctx.Println(err)
		return
	}
	ctx.Println(config)
}

func NewFilePayload(file string, p micro.Payload) (micro.Payload, error) {

	st, err := os.Stat(file)

	if err != nil {
		return nil, err
	}

	config, err := GetConfigWithFile(file)

	if err != nil {
		return nil, err
	}

	err = p.SetConfig(config)

	if err != nil {
		return nil, err
	}

	c := make(chan int8)

	go func() {

		T := time.NewTicker(time.Second * 12)

		loopbreak := false

		for !loopbreak {
			select {
			case <-T.C:
				s, _ := os.Stat(file)
				if s != nil && s.ModTime() != st.ModTime() {
					st = s
					updateConfig(file, p)
				}
			case <-c:
				loopbreak = true
			}
		}

		close(c)

	}()

	return &filePayload{p: p, c: c}, nil
}

func (p *filePayload) SetConfig(config interface{}) error {
	return p.p.SetConfig(config)
}

func (p *filePayload) NewContext(name string, trace string) (micro.Context, error) {
	return p.p.NewContext(name, trace)
}

func (p *filePayload) Exit() {
	p.c <- 1
	p.p.Exit()
}

func (p *filePayload) GetValue(key string) interface{} {
	return p.p.GetValue(key)
}

func (p *filePayload) SetValue(key string, value interface{}) {
	p.p.SetValue(key, value)
}
