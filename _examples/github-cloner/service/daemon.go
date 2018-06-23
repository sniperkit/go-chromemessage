package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// internal - core
	daemon "github.com/sniperkit/snk.golang.daemon/pkg"
)

const (
	defaultServiceName = "chrome-github-cloner"
	defaultServiceDesc = "My chrome backend service for github cloner/sniperkit."
	defaultServiceAddr = ":9977"
)

// serviceDependencies that are NOT required by the service, but might be used
var serviceDependencies = []string{"dummy.service"}

// service has embedded daemon
type Service struct {
	// private
	daemon.Daemon

	// exported
	Config *Config `json:"config" yaml:"config" toml:"config"`
}

type ServiceInfo struct {
	Name string `json:"name" yaml:"name" toml:"name"`
	Desc string `default:""My Github Cloner - Service"" json:"description" yaml:"description" toml:"description"`
}

type DaemonConfig struct {
	Deps string `json:"dependencies" yaml:"dependencies" toml:"dependencies"`
}

type ServiceConfig struct {
	Addr string `default:":9977" json:"address" yaml:"address" toml:"address"`
	host string `json:"-" yaml:"-" toml:"-"`
	port int    `json:"-" yaml:"-" toml:"-"`
}

type Config struct {
	Info    *ServiceInfo   `json:"info" yaml:"info" toml:"info"`
	Daemon  *DaemonConfig  `json:"daemon" yaml:"daemon" toml:"daemon"`
	Service *ServiceConfig `json:"service" yaml:"service" toml:"service"`
}

var defaultConfig = &Config{
	Info: &ServiceInfo{
		Name: defaultServiceName,
		Desc: defaultServiceDesc,
	},
	Daemon: &DaemonConfig{
		Deps: strings.Join(serviceDependencies, ","),
	},
	Service: &ServiceConfig{
		Addr: defaultServiceAddr,
	},
}

func newService(conf *Config) *Service {
	if conf == nil {
		conf = defaultConfig
	}
	dependencies := strings.Split(conf.Daemon.Deps, ",")
	srv, err := daemon.New(conf.Info.Name, conf.Info.Desc, dependencies...)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	svc := &Service{
		srv, &Config{},
	}
	return svc
}

// Manage by daemon commands or run the daemon
func (svc *Service) Manage(command string) (string, error) {
	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", svc.Config.Info.Name)
	// if received any kind of command, do it

	switch command {
	case "install":
		return svc.Install()

	case "remove":
		return svc.Remove()

	case "start":
		return svc.Start()

	case "stop":
		return svc.Stop()

	case "status":
		return svc.Status()

	default:
		return usage, nil

	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined addr, host and port will be extracted later
	listener, err := net.Listen("tcp", svc.Config.Service.Addr)
	if err != nil {
		// ?!!!!!
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)

		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal)
			log.Println("Stoping listening on ", listener.Addr())

			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}
