package core

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type IpcServer struct {
	path string
}

func NewIpcServer(socketPath string) *IpcServer {
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer listener.Close()
		if err := http.Serve(listener, nil); err != nil {
			log.Panicln(err)
		}
	}()

	s := &IpcServer{path: socketPath}
	return s
}

type IpcClient struct {
	Client *http.Client
	path   string
}

func NewIpcClient(socketPath string) *IpcClient {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	c := &IpcClient{path: socketPath, Client: client}
	return c
}
