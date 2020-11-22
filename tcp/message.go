package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Request struct {
	service string
	data    Packable
}

func (req *Request) String() string {
	return fmt.Sprintf("Request{service: %v, data: %v}", req.service, req.data)
}

type Response struct {
	status int32
	data   Packable
}

func (resp *Response) String() string {
	return fmt.Sprintf("Response{status: %v, data: %v}", resp.status, resp.data)
}

type Packable interface {
	Pack() []byte
	UnPack(data []byte)
}

func (req *Request) Pack() []byte {

	buffer := bytes.NewBuffer([]byte{})

	var lth = int32(len(req.service))
	_ = binary.Write(buffer, binary.BigEndian, &lth)
	_ = binary.Write(buffer, binary.BigEndian, []byte(req.service))

	data := req.data.Pack()
	lth = int32(len(data))
	_ = binary.Write(buffer, binary.BigEndian, &lth)
	_ = binary.Write(buffer, binary.BigEndian, data)

	return buffer.Bytes()
}

func (req *Request) UnPack(reqBytes []byte) {
	var lth int32

	reader := bytes.NewReader(reqBytes)

	_ = binary.Read(reader, binary.BigEndian, &lth)
	service := make([]byte, int(lth))
	_ = binary.Read(reader, binary.BigEndian, &service)

	_ = binary.Read(reader, binary.BigEndian, &lth)
	data := make([]byte, int(lth))
	_ = binary.Read(reader, binary.BigEndian, &data)

	req.service = string(service)
	if nil != req.data {
		req.data.UnPack(data)
	}
}

func (resp *Response) Pack() []byte {

	buffer := bytes.NewBuffer([]byte{})

	_ = binary.Write(buffer, binary.BigEndian, &resp.status)

	data := resp.data.Pack()
	lth := int32(len(data))
	_ = binary.Write(buffer, binary.BigEndian, &lth)
	_ = binary.Write(buffer, binary.BigEndian, data)

	return buffer.Bytes()
}

func (resp *Response) UnPack(respBytes []byte) {

	reader := bytes.NewReader(respBytes)

	var status int32
	_ = binary.Read(reader, binary.BigEndian, &status)

	var lth int32
	_ = binary.Read(reader, binary.BigEndian, &lth)
	data := make([]byte, int(lth))
	_ = binary.Read(reader, binary.BigEndian, &data)

	resp.status = status
	if nil != resp.data {
		resp.data.UnPack(data)
	}
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
