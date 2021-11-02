package login

type loginStruct struct {
	Email string
}

type verification struct {
	EncryptedCode string
	DecryptedCode string
	SentDate      string
}
