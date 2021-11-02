package login

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"loginWithEmailGo/helpers"
	"net/smtp"
	"time"
)

func LoginSetup(s *gin.RouterGroup)  {
	s.POST("/sendCode",sendCode)
	s.POST("/checkCode",checkVerification)
}

//u can define any key with 32 bits
var key = "alisultaniseviyor"

func sendCode(c *gin.Context)  {
	body :=loginStruct{}
	data,err := c.GetRawData()
	if err != nil {
		helpers.MyAbort(c,"Something went wrong when you are getting values from body")
		return
	}

	err = json.Unmarshal(data,&body)
	if err != nil {
		helpers.MyAbort(c, "Bad Input")
		return
	}

	// cipher key it means you can define 32 digit
	key := "thisis32bitlongpassphraseimusing"

	if !helpers.EmailIsValid(body.Email) {
		if err != nil {
			helpers.MyAbort(c, "Please check your email")
			return
		}
	}

	// getting current time
	currentTime := time.Now().Format("2006-01-02 3:4:5 PM")

	//generating code
	code, _ := helpers.GenerateDigit(6)

	//sending the code to email
	sendEmail(code, body.Email)

	//generating encrypted code from the code
	encryptedCode, err := helpers.EncryptAES([]byte(code), []byte(key))
	if err != nil {
		helpers.MyAbort(c, "Something went wrong when encrypting code")
		return
	}

	//sent the current time as well to save on local storage or phone storage.
	//You can save the time on ur database to check it
	c.JSON(200, gin.H{
		"encryptedCode": encryptedCode,
		"sent_date":     currentTime,
	})
}

func sendEmail(code, mail string) {
	//put ur e-mail address that you want to sent e-mail by.
	from := "emailverifiy8@gmail.com"
	pass := "emailonayla"

	to := []string{
		mail,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := []byte("To: " + mail + "\r\n" +
		"Subject: Verification Code\r\n" +
		"\r\n" +
		"Hello dear,\r\n" + "Your code is\n" +
		code)

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Is Successfully sent.")

}

func checkVerification(c *gin.Context)  {
	body := verification{}
	data,err := c.GetRawData()
	if err != nil {
		helpers.MyAbort(c,"Input format is wrong")
		return
	}
	err= json.Unmarshal(data,&body)
	if err != nil {
		helpers.MyAbort(c,"Bad Format")
		return
	}

	date := "2006-01-02 3:4:5 PM"
	currentTime := time.Now().Format(date)

	// to compare two times according to same Location
	sentDate, err := time.Parse(date, body.SentDate)
	currentTimeParse, err := time.Parse(date, currentTime)

	//getting the differences
	diff := currentTimeParse.Sub(sentDate)

	//getting differences as seconds
	second := int(diff.Seconds())

	encryptCode, _ := base64.StdEncoding.DecodeString(body.EncryptedCode)

	decryptedCode, err := helpers.DecryptAES(encryptCode, []byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	code := body.DecryptedCode

	if second > 30 {
		c.JSON(400, "Your code is expired")
		return
	} else {
		if code == string(decryptedCode) {
			c.JSON(200, "Verification is completed!")
		} else {
			c.JSON(400, "Check your code !!")
		}
	}
}











