package client

import "C"

import (
	"github.com/wendy512/iec61850"
	"sync"
)

type Client struct {
	conn    *iec61850.Connection
	daCache sync.Map
}

func NewClient(settings *iec61850.Settings) (*Client, error) {
	conn, err := iec61850.NewConnection(settings)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn: conn,
	}
	return c, nil
}

func (c *Client) Write(objectReference string, fc iec61850.FC, value interface{}) error {
	return c.conn.Write(objectReference, fc, value)
}

func (c *Client) Read(objectReference string, fc iec61850.FC) (interface{}, error) {
	return c.conn.Read(objectReference, fc)
}

func (c *Client) Close() {
	c.conn.Close()
}
