package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

type Packer interface {
	Pack() []byte
}

type UnPacker interface {
	UnPack(data []byte)
}

type User struct {
	id         int64
	name       string
	password   string
	createDate time.Time
	updateDate time.Time
}

func (user *User) String() string {
	return fmt.Sprintf("User{id: %v, name: %v, password: %v, createDate: %v, udpateDate: %v}", user.id, user.name, user.password, user.createDate, user.updateDate)
}

func (user *User) Pack() []byte {

	buffer := bytes.NewBuffer([]byte{})

	_ = binary.Write(buffer, binary.BigEndian, &user.id)

	var lth = int32(len(user.name))
	_ = binary.Write(buffer, binary.BigEndian, &lth)
	_ = binary.Write(buffer, binary.BigEndian, []byte(user.name))

	lth = int32(len(user.password))
	_ = binary.Write(buffer, binary.BigEndian, &lth)
	_ = binary.Write(buffer, binary.BigEndian, []byte(user.password))

	cd := user.createDate.Unix()
	ud := user.updateDate.Unix()
	_ = binary.Write(buffer, binary.BigEndian, &cd)
	_ = binary.Write(buffer, binary.BigEndian, &ud)

	return buffer.Bytes()
}

func (user *User) UnPack(data []byte) {

	reader := bytes.NewReader(data)

	_ = binary.Read(reader, binary.BigEndian, &user.id)

	var lth int32
	_ = binary.Read(reader, binary.BigEndian, &lth)
	name := make([]byte, int(lth))
	_ = binary.Read(reader, binary.BigEndian, &name)

	_ = binary.Read(reader, binary.BigEndian, &lth)
	password := make([]byte, int(lth))
	_ = binary.Read(reader, binary.BigEndian, &password)

	var cd int64
	_ = binary.Read(reader, binary.BigEndian, &cd)

	var ud int64
	_ = binary.Read(reader, binary.BigEndian, &ud)

	user.name = string(name)
	user.password = string(password)
	user.createDate = time.Unix(cd, 0)
	user.updateDate = time.Unix(ud, 0)

}

func TestUser(t *testing.T) {
	user := &User{1, "2", "3", time.Now(), time.Now()}
	log.Println(user)
	data := user.Pack()
	log.Print(data)

	user2 := &User{}
	user2.UnPack(data)
	log.Println(user2)
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

				user := &User{int64(i), "2", "3", time.Now(), time.Now()}

				data := user.Pack()
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
			user.UnPack(data)
			log.Printf("message [%v] from server %s", user, addr)
		}
	} else {
		log.Fatal(err)
	}
}
