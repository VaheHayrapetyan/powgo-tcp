package powgo_tcp

import (
	"errors"
	pow "github.com/VaheHayrapetyan/powgo"
	"github.com/kataras/golog"
	"net"
)

func PowClient(host string, serverData []byte) (string, error) {
	connection, err := net.Dial("tcp", host)
	if err != nil {
		return "", err
	}
	defer connection.Close()
	conn := Conn{connection: connection}

	golog.Infof("POW CLIENT: CLIENT INFO: %s %s", "client connected to server on host", host)

	//Ping
	err = conn.Write(PingC, []byte(""))
	if err != nil {
		return "", err
	}
	golog.Infof("POW CLIENT: COMMAND INFO: %s %s` %d %s", "client sent command", commandState[PingC], PingC, "to server")

	//Challenge
	command, challenge, err := conn.Read()
	golog.Infof("POW CLIENT: COMMAND INFO: %s %s` %d", "server sent command", commandState[command], command)

	if command == ErrorC {
		return "", errors.New(string(challenge))
	}
	if command != ChallengeC {
		return "", errors.New("server wrong command")
	}

	//Solve
	proof, err := pow.Solve(challenge, serverData)
	if err != nil {
		return "", err
	}

	//Proof
	err = conn.Write(ProofC, proof)
	if err != nil {
		return "", err
	}
	golog.Infof("POW CLIENT: COMMAND INFO: %s %s` %d %s", "client sent command", commandState[ProofC], ProofC, "to server")

	//Pong
	command, response, err := conn.Read()
	golog.Infof("POW CLIENT: COMMAND INFO: %s %s` %d", "server sent command", commandState[command], command)

	if command == ErrorC {
		return "", errors.New(string(response))
	}
	if command != PongC {
		return "", errors.New("server wrong command")
	}

	return string(response), nil
}
