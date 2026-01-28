// SPDX-License-Identifier: GPL-3.0-or-later

package memcached

import (
	"bytes"
	"os"
	"strings"

	"fmt"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/socket"
	"time"
)

type memcachedConn interface {
	connect() error
	disconnect()
	queryStats() ([]byte, error)
}

func newMemcachedConn(conf Config) memcachedConn {
	return &memcachedClient{conn: socket.New(socket.Config{
		Address: conf.Address,
		Timeout: conf.Timeout.Duration(),
	})}
}

type memcachedClient struct {
	conn socket.Client
	last time.Time
}

func (c *memcachedClient) connect() error {
	return c.conn.Connect()
}

func (c *memcachedClient) disconnect() {
	_ = c.conn.Disconnect()
}

func (c *memcachedClient) queryStats() ([]byte, error) {
	now := time.Now()

	if !c.last.IsZero() {
		fmt.Fprintf(os.Stderr, "[memcached][queryStats] now=%s delta=%s\n",
			now.Format(time.RFC3339Nano),
			now.Sub(c.last),
		)
	} else {
		fmt.Fprintf(os.Stderr, "[memcached][queryStats] now=%s first_call\n",
			now.Format(time.RFC3339Nano),
		)
	}
	c.last = now

	var b bytes.Buffer
	if err := c.conn.Command("stats\r\n", func(bytes []byte) (bool, error) {
		s := strings.TrimSpace(string(bytes))
		b.WriteString(s)
		b.WriteByte('\n')

		return !(strings.HasPrefix(s, "END") || strings.HasPrefix(s, "ERROR")), nil
	}); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
