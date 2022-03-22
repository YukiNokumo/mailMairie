package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

//go:embed mail.txt
var Message string

//go:embed mailList.txt
var MailList string

type Mairie struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties struct {
		Email string `json:"email"`
	} `json:"properties"`
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

func main() {
	var deps []string = []string{"18", "36", "23", "87", "19", "03", "63", "15", "46", "12"}
	var ma Mairie = Mairie{}
	var f *os.File
	var err error
	var countTotal int = 0
	var count int = 0

	err = godotenv.Load(".env")

	if f, err = os.Create("./mailList.txt"); err != nil {
		os.Exit(2)
	}

	defer f.Close()

out:
	for _, dep := range deps {
		if resp, err := http.Get("https://etablissements-publics.api.gouv.fr/v3/departements/" + dep + "/mairie"); err != nil {
			fmt.Println(err)
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			json.Unmarshal(body, &ma)
			for _, v := range ma.Features {
				countTotal++
				if strings.Contains(MailList, v.Properties.Email) || strings.Contains(v.Properties.Email, "http") {
					//fmt.Println(v.Properties.Email + " exist. already send")
					continue
				}
				if err := send(Message, v.Properties.Email); err != nil {
					fmt.Println(err)
					break out
				}
				time.Sleep(1 * time.Second)
				count++
				_, _ = f.Write([]byte(v.Properties.Email + "\n"))
			}
		}
	}

	_, _ = f.Write([]byte(MailList + "\n"))

	f.Sync()

	fmt.Println(count)
	fmt.Println(countTotal)
	fmt.Println("Finished")
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func send(mess string, to1 string) error {
	// Sender data.
	from := os.Getenv("EMAILFROM") //example
	password := os.Getenv("PASSWORDSMTP")
	// Receiver email address.
	to := []string{
		to1,
	}
	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Message.
	message := []byte("From: xxxxx\r\n" +
		"To: " + to1 + "\r\n" +
		"Subject: Recherche Terrain\r\n\r\n" +
		mess + "\r\n")
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent To : " + to1)

	return nil
}
