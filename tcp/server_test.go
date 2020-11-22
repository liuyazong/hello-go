package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	if listener, err := net.Listen("tcp", "127.0.0.1:8888"); err == nil {
		serverAddr := listener.Addr().String()
		defer func() {
			_ = listener.Close()
			log.Printf("Server closed: %s", serverAddr)
		}()
		log.Printf("Server started: %s", serverAddr)
		for {
			conn, _ := listener.Accept()
			addr := conn.RemoteAddr().String()
			go func(conn net.Conn) {
				defer func() {
					_ = conn.Close()
					log.Printf("%s closed", addr)
				}()
				log.Printf("%s connected", addr)

				scanner := bufio.NewScanner(conn)

				scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
					if !atEOF {
						if len(data) > 4 {
							var lth int32
							_ = binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &lth)
							if int(lth) <= len(data) {
								return int(lth), data[4:int(lth)], nil
							}
						}
					} else {
						log.Printf("eof: %v", atEOF)
					}
					return
				})

				for buffer := bytes.NewBuffer([]byte{}); scanner.Scan(); buffer.Reset() {

					data := scanner.Bytes()

					user := &User{}
					req := &Request{data: user}
					req.UnPack(data)
					log.Printf("message [%v] from client %s", req, addr)

					user.id = user.id + 1

					resp := &Response{status: 1, data: user}
					data = resp.Pack()

					lth := int32(len(data) + 4)

					_ = binary.Write(buffer, binary.BigEndian, &lth)
					_ = binary.Write(buffer, binary.BigEndian, &data)

					data = buffer.Bytes()

					writer := bufio.NewWriter(conn)
					_, _ = writer.Write(data)
					_ = writer.Flush()
				}
			}(conn)
		}
	} else {
		log.Fatal(err)
	}
}
