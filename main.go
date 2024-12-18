package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func bytesToHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}
func PemToDer(pemData []byte) ([]byte, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	return block.Bytes, nil
}

func setPublicKey(dataSize int, encryptedData []byte, encryptionKey []byte) (*rsa.PublicKey, error) {

	// AES decryption using CBC mode and PKCS5 padding
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Creating the cipher and initialization vector (IV)
	blockMode := cipher.NewCBCDecrypter(block, encryptionKey)
	decryptedData := make([]byte, dataSize)

	// Decrypting data
	blockMode.CryptBlocks(decryptedData, encryptedData)

	// Extract the decrypted data (excluding padding)
	decryptedDataWithoutPadding := decryptedData[:len(encryptedData)]
	decryptedDataString := string(decryptedDataWithoutPadding)
	fmt.Println("Decrypted data without padding:", decryptedDataString)

	derData, err := PemToDer(decryptedDataWithoutPadding)
	if err != nil {
		return nil, err
	}

	// Parse the DER data to get the RSA public key
	rsakey, err := x509.ParsePKIXPublicKey(derData)
	if err != nil {
		return nil, err
	}

	// Assert the type of the public key as RSA
	publicKey, ok := rsakey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return publicKey, nil
}

func main() {
	// Define the address of the TCP server
	serverAddr := "10.54.2.161:53150" // Change this to your server's address
	// Connect to the TCP server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddr)
	var3 := time.Now().UnixMilli()

	// Convert to big-endian (network byte order)
	var5 := make([]byte, 8)
	binary.BigEndian.PutUint64(var5, uint64(var3))

	// Create a SessionID array of size 16
	SessionID := make([]byte, 16)

	// Copy var5 into the first and second 8 bytes of SessionID
	copy(SessionID[0:8], var5) // First 8 bytes
	copy(SessionID[8:16], var5)

	SessionIDString := bytesToHex(SessionID)
	hexString := "41470102000000000008000000080000000100000001280000000000ed5793750100080000000300"
	hexString2 := "41470102000100000016000000160000000100000001280200000000ed4a937d090016000000"
	//       _ := "41470102000100000016000000160000000100000001280200000000ED4A937D0900160000003E6839DA930100003E6839DA93010000"

	payload2 := hexString2 + SessionIDString

	// Convert the hex string to a byte slice
	rawRequest, err := hex.DecodeString(hexString)
	if err != nil {
		log.Fatalf("Error decoding hex string: %v\n", err)
	}
	rawRequest2, err := hex.DecodeString(payload2)
	if err != nil {
		log.Fatalf("Error decoding hex string: %v\n", err)
	}

	// Send the raw byte request
	_, err = conn.Write(rawRequest)
	if err != nil {
		log.Fatalf("Error sending data: %v\n", err)
	}

	fmt.Println("Raw byte request sent")

	// Optionally, read the response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading response: %v\n", err)
	}

	// Print the server's response
	fmt.Printf("Server response: %s\n", string(buffer[:n]))

	_, err = conn.Write(rawRequest2)
	if err != nil {
		log.Fatalf("Error sending data: %v\n", err)
	}

	fmt.Println("Raw byte request sent")

	buffer2 := make([]byte, 1024)
	n2, err := conn.Read(buffer2)
	if err != nil {
		log.Fatalf("Error reading response: %v\n", err)
	}

	realbuffer := buffer2[70:n2]
	key := []byte(SessionID)

	rsakey, err := setPublicKey(n2-6, realbuffer, key)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Public Key:", rsakey)
	}
}
