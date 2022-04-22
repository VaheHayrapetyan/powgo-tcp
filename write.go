package powgo_tcp

import (
	"encoding/binary"
)

// Write with command and message
func (conn *Conn) Write(command Command, message []byte) error {
	protoMess := protocolFormat(command, message) // string pars to protocol format
	return conn.write(protoMess)                  // writing mess
}

// All Message Writing Without Loss
func (conn *Conn) write(mess []byte) error {
	start := 0
	finish := len(mess)

	for {
		sentBytesCount, err := conn.connection.Write(mess[start:]) // writing process
		if err != nil {
			return err
		}

		start += sentBytesCount
		if start == finish { //all message wrote
			break
		}
	}

	return nil
}

//protocolFormat string(message) parse to Protocol Type
func protocolFormat(command Command, mess []byte) []byte {
	lengthUInt32 := uint32(len(mess)) //message length
	var lengthByte = make([]byte, 4)
	binary.BigEndian.PutUint32(lengthByte, lengthUInt32) // message length parsing

	return []byte(string(command) + string(lengthByte) + string(mess))
}
