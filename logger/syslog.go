package logger

import (
	"fmt"
	"log"
	"log/syslog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/matt-abi/abi-lib/dynamic"
)

type syslogConfig struct {
	Net  string `json:"net"`
	Addr string `json:"addr"`
	Tag  string `json:"tag"`
	Dns  string `json:"dns"`
}

type syslogLogger struct {
	E chan bool
	C chan *string
	W *sync.WaitGroup
}

func NewSyslogLogger(config interface{}) (Logger, error) {

	E := make(chan bool)
	C := make(chan *string, 2048)
	W := &sync.WaitGroup{}

	cfg := &syslogConfig{}
	dynamic.SetValue(cfg, config)

	if cfg.Tag == "" {
		cfg.Tag = "-"
	}

	lock := sync.RWMutex{}
	addrs := map[string]bool{}

	leave := func(addr string) {
		lock.Lock()
		delete(addrs, addr)
		lock.Unlock()
	}

	join := func(addr string) {
		lock.Lock()
		b := addrs[addr]
		if !b {
			addrs[addr] = true
		}
		lock.Unlock()
		if !b {

			W.Add(1)

			go func() {

				defer W.Done()

				var w *syslog.Writer = nil
				var err error = nil
				tk := time.NewTicker(6 * time.Second)
				defer tk.Stop()

				running := true

				for running {

					if w == nil {
						lock.RLock()
						b := addrs[addr]
						lock.RUnlock()
						if !b {
							running = false
							break
						}
						w, err = syslog.Dial(cfg.Net,
							addr,
							syslog.LOG_INFO|syslog.LOG_WARNING|syslog.LOG_DEBUG|syslog.LOG_NOTICE|syslog.LOG_ERR,
							cfg.Tag)
						if err != nil {
							log.Printf("syslog error %s %s wait 6 seconds to try again", addr, err.Error())
						}
					}

					if w == nil {
						select {
						case <-tk.C:
						case _, ok := <-E:
							if !ok {
								running = false
								break
							}
						}
					} else {
						select {
						case _, ok := <-E:
							if !ok {
								running = false
								break
							}
						case s := <-C:
							if s == nil {
								running = false
								break
							} else {
								_, err = fmt.Fprintln(w, *s)
								if err != nil {
									w.Close()
									w = nil
									log.Printf("syslog write error %s %s", addr, err.Error())
								}
							}
						}
					}
				}

				if w != nil {
					w.Close()
				}
			}()
		}
	}

	if cfg.Dns != "" {

		ps := strings.Split(cfg.Dns, ":")
		pn := len(ps)

		W.Add(1)

		go func() {

			defer W.Done()

			tk := time.NewTicker(6 * time.Second)
			defer tk.Stop()

			vs := map[string]bool{}

			running := true

			for running {

				ns, err := net.LookupIP(ps[0])

				if err != nil {
					log.Printf("dns lookup error %s %s wait 6 seconds to try again", ps[0], err.Error())
				} else {
					s := map[string]bool{}
					for _, n := range ns {
						ip := n.String()
						s[ip] = true
						vs[ip] = true
						if pn > 1 {
							join(fmt.Sprintf("%s:%s", n, ps[1]))
						} else {
							join(fmt.Sprintf("%s:514", n))
						}
					}
					ks := []string{}
					for v := range vs {
						if !s[v] {
							leave(v)
							ks = append(ks, v)
						}
					}
					for _, v := range ks {
						delete(vs, v)
					}
				}

				select {
				case <-tk.C:
				case _, ok := <-E:
					if !ok {
						running = false
					}
				}

			}

		}()

	}

	for _, v := range strings.Split(cfg.Addr, ",") {
		if v != "" {
			join(v)
		}
	}

	return &syslogLogger{C: C, W: W, E: E}, nil
}

func (l *syslogLogger) Output(text string) {
	l.C <- &text
}

func (l *syslogLogger) Recycle() {
	log.Println("SyslogLogger Recycle ...")
	if l.E != nil {
		close(l.E)
		l.E = nil
		l.W.Wait()
		l.W = nil
		close(l.C)
		l.C = nil
	}
	log.Println("SyslogLogger Recycle done")
}
