package packet

type EncryptionResponse struct {
	SharedSecret []byte `proto:"varIntByteArray"`
	VerifyToken  []byte `proto:"varIntByteArray"`
}

func (er *EncryptionResponse) GetPacketId() int {
	return 0x01
}

/**
func (er *EncryptionResponse) Handle(packet []byte, connection *Connection) {
	fmt.Println("Shared Secret ", er.SharedSecret)
	fmt.Println("Verify Token", er.VerifyToken)

	decryptedSecret, exception := rsa.DecryptPKCS1v15(rand.Reader, connection.Key, er.SharedSecret)

	connection.SharedSecret = er.SharedSecret
	if exception != nil {
		fmt.Println("Error decrypting secret token:", exception)
		return
	}

	fmt.Println("Decrypted Secret ", decryptedSecret)

	decryptedToken, exception := rsa.DecryptPKCS1v15(rand.Reader, connection.Key, er.VerifyToken)

	if exception != nil {
		fmt.Println("Error decrypting verify token:", exception)
		return
	}

	fmt.Println("Decrypted Token ", decryptedToken)

	for i, val := range connection.VerifyToken {
		if val != decryptedToken[i] {
			fmt.Println("Token is not valid!!")
			return
		}
	}
	fmt.Println("Encryption is verified!")

	publicKey, exception := x509.MarshalPKIXPublicKey(&connection.Key.PublicKey)

	if exception != nil {
		fmt.Println("Marshalling public key", exception)
		return
	}

	hash := sha1.New()
	hash.Write([]byte(""))
	hash.Write(decryptedSecret)
	hash.Write(publicKey)
	sum := hash.Sum(nil)
	negative := (sum[0] & 0x80) == 0x80
	if negative {
		sum = twosComplement(sum)
	}

	res := strings.TrimLeft(fmt.Sprintf("%x", sum), "0")
	if negative {
		res = "-" + res
	}

	fmt.Println("Authenticating session", res)

	sessionUrl := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", connection.Player.Username, res)
	fmt.Println(sessionUrl)
	response, exception := http.Get(sessionUrl)

	if exception != nil {
		fmt.Println("Exception contacting session server ", exception)
		return
	}

	body := make([]byte, response.ContentLength)

	response.Body.Read(body)

	sessionResponse := SessionResponse{}

	exception = json.Unmarshal(body, &sessionResponse)

	if exception != nil {
		fmt.Println("Exception parsing session server response ", exception)
		return
	}

	fmt.Println("Authed as ", sessionResponse.Name)

	uuid, exception := hex.DecodeString(sessionResponse.ID)
	connection.Player.UUID = uuid

	if exception != nil {
		fmt.Println("Exception parsing UUID ", exception)
		return
	}

	connection.Player.Username = sessionResponse.Name
	connection.Player.DisplayName = entity.ChatMessageComponent{
		Text: sessionResponse.Name,
	}

	connection.Player.Properties = sessionResponse.Properties

	aesCipher, exception := aes.NewCipher(decryptedSecret)

	if exception != nil {
		fmt.Println("Exception creating cipher ", exception)
		return
	}

	connection.Cipher = cfb8.NewCFB8Encrypt(aesCipher, decryptedSecret)
	connection.DecryptionCipher = cfb8.NewCFB8Decrypt(aesCipher, decryptedSecret)
	connection.EnableEncryption = true

	connection.Minecraft.StartPlayerJoin(connection)
}

// little endian
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}

type SessionResponse struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	Properties []entity.PlayerProperty `json:"properties"`
}
*/
