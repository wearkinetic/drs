package drs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Server struct {
	*Processor
	Protocol   protocol.Protocol
	Transport  Transport
	connect    []func(*Connection) error
	disconnect []func(*Connection)
	inbound    cmap.ConcurrentMap
}

func New(transport Transport, protocol protocol.Protocol) *Server {
	result := &Server{
		Processor:  NewProcessor(),
		Protocol:   protocol,
		Transport:  transport,
		connect:    []func(*Connection) error{},
		disconnect: []func(*Connection){},
		inbound:    cmap.New(),
	}
	http.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		functions := []string{}
		for key, _ := range result.Processor.handlers {
			functions = append(functions, key)
		}
		host, _ := os.Hostname()
		response(w, 200, map[string]interface{}{
			"hostname": host,
			"connections": map[string]interface{}{
				"inbound": result.inbound.Count(),
			},
			"commands":  result.stats.Items(),
			"functions": functions,
		})
	})
	return result
}

func (this *Server) Listen(host string) error {
	return this.Transport.Listen(host, func(raw io.ReadWriteCloser) {
		conn := Accept(this.Protocol, raw)
		conn.Processor = this.Processor
		for _, cb := range this.connect {
			err := cb(conn)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
		}
		key := uuid.Ascending()
		defer func() {
			this.inbound.Remove(key)
			for _, cb := range this.disconnect {
				cb(conn)
			}
		}()
		this.inbound.Set(key, conn)
		conn.Read()
	})
}

func (this *Server) OnConnect(cb func(*Connection) error) {
	this.connect = append(this.connect, cb)
}

func (this *Server) OnDisconnect(cb func(*Connection)) {
	this.disconnect = append(this.disconnect, cb)
}

func (this *Server) Close() {

}

func (this *Server) Broadcast(msg *Command) {
	for kv := range this.inbound.Iter() {
		kv.Val.(*Connection).Fire(msg)
	}
}

func response(w http.ResponseWriter, status int, input interface{}) {
	data, _ := json.Marshal(input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
