package tcp

import (
	"errors"
)

// Set of error variables for start up.
var (
	ErrInvalidNetType       = errors.New("invalid net type configuration")
	ErrInvalidConnHandler   = errors.New("invalid connection handler configuration")
	ErrInvalidReqHandler    = errors.New("invalid request handler configuration")
	ErrInvalidRespHandler   = errors.New("invalid response handler configuration")
	ErrInvalidLoggerHandler = errors.New("invalid logger handler configuration")
)

// Logger defines an handler used to help log events.
type Logger func(evt string, typ string, ipAddress string, format string, a ...any)

// Config provides a data structure of required configuration parameters.
type Config struct {
	NetType     string      // "tcp", tcp4" or "tcp6"
	Addr        string      // "host:port" or "[ipv6-host%zone]:port"
	ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
	ReqHandler  ReqHandler  // Support for handling the specific request workflow.
	RespHandler RespHandler // Support for handling the specific response workflow.
	Logger      Logger      // Support for logging events that occur in the TCP listener.
}

func (cfg Config) validate() error {
	if cfg.NetType != "tcp" && cfg.NetType != "tcp4" && cfg.NetType != "tcp6" {
		return ErrInvalidNetType
	}

	if cfg.ConnHandler == nil {
		return ErrInvalidConnHandler
	}

	if cfg.ReqHandler == nil {
		return ErrInvalidReqHandler
	}

	if cfg.RespHandler == nil {
		return ErrInvalidRespHandler
	}

	if cfg.Logger == nil {
		return ErrInvalidLoggerHandler
	}

	return nil
}
