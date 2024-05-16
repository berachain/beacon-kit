package client

import (
	"context"
	"net"
	"net/rpc"
	"os"
)

//nolint:lll // long line length due to struct tags.
func (s *EngineClient[ExecutionPayloadDenebT]) startIPCServer(ctx context.Context) {
	// remove existing socket file if exists
	// alternatively we can use existing one by checking for os.IsNotExist(err)
	if _, err := os.Stat(s.cfg.IPCPath); err != nil {
		s.logger.Info("Removing existing IPC file", "path", s.cfg.IPCPath)
		os.Remove(s.cfg.IPCPath)
	}

	// use UDS for IPC
	listener, err := net.Listen("unix", s.cfg.IPCPath)
	if err != nil {
		s.logger.Error("failed to listen on IPC socket", "err", err)
		return
	}
	s.ipcListener = listener

	// register the RPC server
	server := rpc.NewServer()
	if err = server.Register(s); err != nil {
		s.logger.Error("failed to register RPC server", "err", err)
		return
	}
	s.logger.Info("IPC server started", "path", s.cfg.IPCPath)

	// start server in a goroutine
	go func() {
		for {
			// continuously accept incoming connections until context is cancelled
			select {
			case <-ctx.Done():
				s.logger.Info("shutting down IPC server")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					s.logger.Error("failed to accept IPC connection", "err", err)
					continue
				}
				go server.ServeConn(conn)
			}
		}
	}()
}
