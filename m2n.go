package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"golang.org/x/net/proxy"
)

func main() {
	from := "anonymous"
	to := "mail2news@dizum.com"
	host := "smtp.dizum.com"
	port := ":2525"
	proxyAddr := "127.0.0.1:9150"

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Usage: m2n < message.txt (with all Usenet headers)")
		os.Exit(1)
	}

	message, _ := ioutil.ReadAll(os.Stdin)

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating SOCKS5 dialer:", err)
		os.Exit(1)
	}

	smtpClient, err := dialSMTP(dialer, host, port)
	if err != nil {
		fmt.Println("Error creating SMTP client:", err)
		os.Exit(1)
	}
	defer smtpClient.Quit()

	if err = smtpClient.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}); err != nil {
		fmt.Println("Error starting TLS:", err)
		os.Exit(1)
	}

	if err = smtpClient.Mail(from); err != nil {
		fmt.Println("Error setting sender:", err)
		os.Exit(1)
	}
	if err = smtpClient.Rcpt(to); err != nil {
		fmt.Println("Error setting recipient:", err)
		os.Exit(1)
	}

	w, err := smtpClient.Data()
	if err != nil {
		fmt.Println("Error preparing data:", err)
		os.Exit(1)
	}
	_, err = w.Write(message)
	if err != nil {
		fmt.Println("Error writing message:", err)
		os.Exit(1)
	}
	err = w.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		os.Exit(1)
	}

	fmt.Println("Message successfully sent!")
}

func dialSMTP(dialer proxy.Dialer, host, port string) (*smtp.Client, error) {
	conn, err := dialer.Dial("tcp", host+port)
	if err != nil {
		return nil, err
	}
	return smtp.NewClient(conn, host)
}