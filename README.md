# powgo-tcp

Hello. Here you will get acquainted with the powgo-tcp library. It is a client-server system protected by Proof of work (PoW) authorization with a challenge-response protocol and running on TCP.

Remark:
This library uses the github.com/VaheHayrapetyan/powgo library, specifically the powgo.NewChallenge, powgo.Solve, powgo.Verify functions to create a challenge, solve a challenge, and verify a proof.

Functions:

    func PowServer(host string, difficulty uint32, serverData []byte, getResponse func() string) error
    func PowClient(host string, serverData []byte) (string, error)

Types:

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

Now let's start from the beginning․

    func PowServer(host string, difficulty uint32, serverData []byte, getResponse func() string) error

The PowServer function gets host, difficulty, serverData, getResponse. Host is a server host. Difficulty is the degree of complexity that determined how difficult it will be for a prover (client) to find the answer to a given challenge. The data is a server (verifier) standard information known to the client (prover). It is used to solving and verifying a challenge. getResponse is a function that returns a string value that is the final response returned to the client.

    func PowClient(host string, serverData []byte) (string, error)

The PowClient function gets a host, which is the server host, data that is defined in the same way as it was defined for the PowServer function.

As you can see, in the types section, there is a Command type and five types of commands:
0. PingC
1. ChallengeC
2. ProofC
3. PongC
4. ErrorC

Describe the progress of the program։
0. the client connects to the server, then sends it the PingC command,
1. the server receives the PingC command and creates a challenge via the powgo.NewChallenge function and sends it to the client along with the ChallengeC command,
2. upon receiving a challenge, the client solves it using the powgo.Solve function, obtains a proof, and sends it to the server along with the ProofC command,
3. the server receives the proof and verifies it using the powgo.Verify function: when the powgo.Verify function returns true, the server sends the expected response to the client along with the PongC command,
4. when the powgo.Verify function is returned false, the server sends an error message to the client along with the ErrorC command․

Our library is based on TCP. Therefore, we created a communication protocol that looks like this: 1 byte command, 4 bytes messageLength of type int, a message whose length is equal to messageLength.

That's all !!! Thank you for your attention.