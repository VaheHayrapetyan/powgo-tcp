# powgo

Hello. Here you will get acquainted with the powgo library. A library based on the idea of Proof of work (PoW).     
Proof of work (PoW) is a form of cryptographic proof in which one party (the prover) proves to others (the verifiers) that a certain amount of a specific computational effort has been expended. Verifiers can subsequently confirm this expenditure with minimal effort on their part.

Functions:

    func NewChallenge(difficulty uint32, nonce []byte) (challenge []byte)
    func Solve(challenge []byte, data []byte) (proof []byte, err error)
	func Verify(challenge []byte, proof []byte, data []byte) (ok bool, err error)

Structures:

    type Algorithm string
    
    const (
        Sha111 Algorithm = "sha111"
    )

    type Challenge struct {
        Alg        Algorithm // The requested algorithm
        Difficulty uint32    // The requested difficulty
        Nonce      []byte    // Nonce to diversify the challenge
    }

    type Proof struct {
        buf []byte
    }

Now let's start from the beginning․     
In order for the provers to prove that they are worthy of something, it is necessary to challenge them, for which they will be rewarded.

    func NewChallenge(difficulty uint32, nonce []byte) (challenge []byte)

The NewChallenge function creates and returns this challenge. Difficulty and nonce parameters are passed into the function to add some complexity and diversify. Difficulty is the degree of complexity that determined how difficult it will be for a prover to find the answer to a given challenge. Nonce is information that is sent to the client (prover) each time a challenge is presented. Through it, the challenges are different from each other.

    func Solve(challenge []byte, data []byte) (proof []byte, err error)

The Solve function solves the challenge and returns the answer as proof. Challenge and data parameters are passed to the function. Data is a server (verifier) standard information known to the client (verifier). It is used to solving and verifying a challenge. The function returns an error, for example, when the challenge has an invalid structure or when the algorithm type is not equal to Sha111.     

    func Verify(challenge []byte, proof []byte, data []byte) (ok bool, err error)

The Verify function checks the solution against the Sha111 algorithm, returns true if the correct solution was sent, and false otherwise. The function takes as parameters the challenge, the proof of the given challenge, and the data which is the same data we passed to the Solve function. The function returns an error, for example, when the challenge has an invalid structure or when the algorithm type is not equal to Sha111.         

    //example of using a powgo library

    import (
        "fmt"
        pow "github.com/VaheHayrapetyan/powgo"
    )
    
    func main() {
        //create a proof of work challenge with difficulty 30
        challenge := pow.NewChallenge(30, []byte("some random nonce"))
        fmt.Printf("challenge: %s\n", challenge)
    
        // solve the proof of work
        proof, err := pow.Solve(challenge, []byte("some server data"))
        if err != nil {
            panic(err)
        }
        fmt.Printf("proof: %s\n", proof)
    
        // verify if the proof is correct
        ok, err := pow.Verify(challenge, proof, []byte("some server data"))
        if err != nil {
            panic(err)
        }
        fmt.Printf("verify: %v", ok)
    }

        // Output:
        //		challenge:  sha2bday-5-c29tZSByYW5kb20gbm9uY2U
        //		proof:      AAAAAAAGZZ8AAAAAAA4LyQAAAAAADyI4
        //		check:      true

That's all !!! Thank you for your attention.