package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := &User{1, "2", "3", time.Now().Local(), time.Now().Local()}
	log.Println(user)
	data := user.Pack()

	user = &User{}
	user.UnPack(data)

	req := &Request{service: "echo", data: user}

	pack := req.Pack()
	log.Println(pack)

	req = &Request{data: &User{}}
	req.UnPack(pack)
	log.Println(req)

	resp := &Response{status: 1, data: user}
	pack = resp.Pack()
	log.Println(pack)

	resp = &Response{data: &User{}}
	resp.UnPack(pack)
	log.Println(resp)

}

func TestClient(t *testing.T) {

	if conn, err := net.Dial("tcp", "127.0.0.1:8888"); err == nil {
		addr := conn.RemoteAddr().String()
		defer func() {
			_ = conn.Close()
			log.Printf("%s closed", addr)
		}()

		log.Printf("%s connected", addr)

		for i := 0; i < 10; i++ {
			go func(i int) {

				user := &User{int64(i), "2", "3", time.Now().Local(), time.Now().Local()}
				req := &Request{service: "echo", data: user}

				data := req.Pack()
				lth := int32(len(data) + 4)

				buffer := bytes.NewBuffer([]byte{})

				_ = binary.Write(buffer, binary.BigEndian, &lth)
				_ = binary.Write(buffer, binary.BigEndian, &data)

				data = buffer.Bytes()

				writer := bufio.NewWriter(conn)
				_, _ = writer.Write(data)
				_ = writer.Flush()

			}(i)
		}

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
		for scanner.Scan() {
			data := scanner.Bytes()
			user := &User{}
			resp := &Response{data: user}
			resp.UnPack(data)
			log.Printf("message [%v] from server %s", resp, addr)
		}
	} else {
		log.Fatal(err)
	}
}
