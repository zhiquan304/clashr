package provider

import (
	"context"
	"errors"
	"github.com/Dreamacro/clash/adapters/inbound"
	"github.com/Dreamacro/clash/component/socks5"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/tunnel"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Vehicle Type
const (
	File VehicleType = iota
	HTTP
	Compatible
)

// VehicleType defined
type VehicleType int

func (v VehicleType) String() string {
	switch v {
	case File:
		return "File"
	case HTTP:
		return "HTTP"
	case Compatible:
		return "Compatible"
	default:
		return "Unknown"
	}
}

type Vehicle interface {
	Read() ([]byte, error)
	Path() string
	Type() VehicleType
}

type FileVehicle struct {
	path string
}

func (f *FileVehicle) Type() VehicleType {
	return File
}

func (f *FileVehicle) Path() string {
	return f.path
}

func (f *FileVehicle) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func NewFileVehicle(path string) *FileVehicle {
	return &FileVehicle{path: path}
}

type HTTPVehicle struct {
	url  string
	path string
}

func (h *HTTPVehicle) Type() VehicleType {
	return HTTP
}

func (h *HTTPVehicle) Path() string {
	return h.path
}

func (h *HTTPVehicle) Read() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, h.url, nil)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: func(_ context.Context, network string, address string) (net.Conn, error) {
			if network != "tcp" && network != "tcp4" && network != "tcp6" {
				return nil, errors.New("Unsupported network type " + network)
			}

			client, server := net.Pipe()

			tunnel.Add(inbound.NewSocket(socks5.ParseAddr(address), server, constant.HTTP, constant.TCP))

			return client, nil
		},
	}

	client := http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func NewHTTPVehicle(url string, path string) *HTTPVehicle {
	return &HTTPVehicle{url, path}
}
