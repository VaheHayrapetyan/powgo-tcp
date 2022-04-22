package gopow_tcp

import (
	"fmt"
	pow "github.com/VaheHayrapetyan/powgo"
	"github.com/kataras/golog"
	"net"
	"sync/atomic"
)

var ops uint64

const maxPacketLength = 1000

type Conn struct {
	id         uint64
	connection net.Conn
}

type Command uint8

const (
	PingC = iota
	ChallengeC
	ProofC
	PongC
	ErrorC
)

var commandState = map[Command]string{
	PingC:      "PING",
	ChallengeC: "CHALLENGE",
	ProofC:     "PROOF",
	PongC:      "PONG",
	ErrorC:     "ERROR",
}

func PowServer(host string, difficulty uint32, serverData []byte, getResponse func() string) error {
	// Listen
	listen, err := net.Listen("tcp4", host)
	if err != nil {
		return err
	}
	defer listen.Close()

	golog.Infof("POW SERVER: SERVER INFO: %s %s", "server started on host", host)

	// Accept
	for {
		connection, err := listen.Accept()
		if err != nil {
			return err
		}

		// Atomic
		atomic.AddUint64(&ops, 1)
		id := atomic.LoadUint64(&ops)

		conn := Conn{
			id:         id,
			connection: connection,
		}
		go conn.handler(difficulty, serverData, getResponse())
	}
}

func (conn *Conn) handler(difficulty uint32, serverData []byte, response string) {
	defer conn.connection.Close()
	golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "connected")

	//Ping
	command, _, err := conn.Read()
	if err != nil {
		golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		return
	}

	golog.Infof("POW SERVER: COMMAND INFO: %s %d %s %s` %d", "client", conn.id, "sent command", commandState[command], command)
	if command != PingC {
		errMess := fmt.Sprintf("client command is %d: first command must be Ping` 0", command)
		golog.Errorf("POW SERVER: CLIENT ERROR: %s", errMess)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		err := conn.Write(ErrorC, []byte(errMess))
		if err != nil {
			golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
			return
		}
		return
	}

	//Challenge
	challenge := pow.NewChallenge(difficulty, newNonce())

	err = conn.Write(ChallengeC, challenge)
	if err != nil {
		golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		return
	}

	golog.Infof("POW SERVER: COMMAND INFO: %s %s` %d %s %d", "server sent command", commandState[ChallengeC], ChallengeC, "to client", conn.id)

	//Proof
	command, proof, err := conn.Read()
	if err != nil {
		golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		return
	}

	golog.Infof("POW SERVER: COMMAND INFO: %s %d %s %s` %d", "client", conn.id, "sent command", commandState[command], command)
	if command != ProofC {
		errMess := fmt.Sprintf("client command is %d: second command mut be Proof` 2", command)
		golog.Errorf("POW SERVER: CLIENT ERROR: %s", errMess)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		err := conn.Write(ErrorC, []byte(errMess))
		if err != nil {
			golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
			return
		}
		return
	}

	//Verify
	verified, err := pow.Verify(challenge, proof, serverData)
	if err != nil {
		golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		err := conn.Write(ErrorC, []byte(err.Error()))
		if err != nil {
			golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
			return
		}
		return
	}

	if !verified {
		errMess := fmt.Sprintf("proof is not verified")
		golog.Errorf("POW SERVER: CLIENT ERROR: %s", errMess)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		err := conn.Write(ErrorC, []byte(errMess))
		if err != nil {
			golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
			return
		}
		return
	}

	//Pong
	err = conn.Write(PongC, []byte(response))
	if err != nil {
		golog.Errorf("POW SERVER: CLIENT ERROR: %v", err)
		golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
		return
	}

	golog.Infof("POW SERVER: COMMAND INFO: %s %s` %d %s %d", "server sent command", commandState[PongC], PongC, "to client", conn.id)
	golog.Infof("POW SERVER: CLIENT INFO: %s %d %s", "client", conn.id, "closed")
}
