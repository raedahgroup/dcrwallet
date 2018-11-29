package dcrtxclient

import (
	"strings"
	"sync"

	"github.com/decred/dcrwallet/dcrtxclient/service"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

func init() {
	grpc.EnableTracing = false
}

type (
	Config struct {
		Enable  bool
		Address string
		Timeout uint32
	}

	Client struct {
		sync.Mutex
		Cfg        *Config
		conn       *grpc.ClientConn
		TxService  *service.TxService
		Ws         *websocket.Conn
		IsShutdown bool
		TxHashes   []string
	}
)

// ContainTxHash checks whether transaction hashes contains the given hash.
func (c *Client) ContainTxHash(txHash string) bool {
	for _, hash := range c.TxHashes {
		if strings.Compare(hash, txHash) == 0 {
			return true
		}
	}
	return false
}

// startSession establishes a connection to the transaction matching server if
// the client's configuration allows it.
func (c *Client) StartSession() error {
	if c.Cfg.Enable {
		conn, err := c.Connect()
		if err != nil {
			return err
		}

		if conn == nil {
			return ErrCannotConnect
		}

		c.conn = conn
		err = c.registerServices()
		if err != nil {
			return err
		}
	}

	if !c.Cfg.Enable {
		return ErrConfigDisable
	}

	return nil
}

// Connect attempts to connect to dcrtxmatcher server
func (c *Client) Connect() (*grpc.ClientConn, error) {
	c.Lock()
	defer c.Unlock()

	conn, err := grpc.Dial(c.Cfg.Address, grpc.WithInsecure())
	if err != nil {
		log.Warn("Unable to connect to dcrtxmatcher server.")
		return nil, err
	}

	return conn, nil
}

// Disconnect disconnects client from server
// returns error if client is not connected
func (c *Client) Disconnect() error {
	if c.isConnected() {
		c.conn.Close()
		log.Info("DcrTxClient grpc disconnected")
	}
	if c.Ws != nil {
		c.Ws.Close()
		log.Info("DcrTxClient websocket disconnected")
	}
	return nil
}

// isConnected checks if client is connected to server
// returns appropriate boolen depending on result
func (c *Client) isConnected() bool {
	if c.conn != nil {
		return true
	}

	return false
}

// registerServices registers service api function with dcrtxmatcher server.
func (c *Client) registerServices() error {
	if !c.isConnected() {
		return ErrNotConnected
	}

	c.TxService = service.NewTxService(c.conn)

	return nil
}
