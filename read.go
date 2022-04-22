package powgo_tcp

import (
	"bufio"
	"encoding/binary"
	"errors"
)

// Read Packet-packet
func (conn *Conn) Read() (Command, []byte, error) {
	reader := bufio.NewReader(conn.connection) //new reader

	commandByte, err := reader.ReadByte() //command type
	if err != nil {
		return ErrorC, nil, err
	}
	command := Command(commandByte) //parse

	packetLengthBuff := make([]byte, 4)
	for i := 0; i < 4; i++ {
		packetByte, err := reader.ReadByte()
		if err != nil {
			return command, nil, err
		}
		packetLengthBuff[i] = packetByte
	}
	packetLength := int(binary.BigEndian.Uint32(packetLengthBuff)) //parse
	if packetLength > maxPacketLength {                            //length error
		return command, nil, errors.New("big packet length")
	}

	packet := make([]byte, packetLength) // declare packet
	for i := 0; i < packetLength; i++ {
		readByte, err := reader.ReadByte()
		if err != nil {
			return command, nil, err
		}

		packet[i] = readByte
	}

	return command, packet[:packetLength], nil
}
