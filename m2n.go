package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "net/smtp"
)

func main() {
    // contacting server as
    from := "user@localhost"

    // sending email to
    to := "mail2news@dizum.com"

    // server we send email through
    host := "smtp.dizum.com"

    // check if there is piped input
    fi, _ := os.Stdin.Stat()
    if (fi.Mode() & os.ModeCharDevice) != 0 {
        fmt.Println("Usage: m2n < message.txt") // with  all Usenet message Headers
        os.Exit(1)
    }

    message, _ := ioutil.ReadAll(os.Stdin)

    if err := smtp.SendMail(host+":25", nil, from, []string{to}, message); err != nil {
        fmt.Println("Error SendMail: ", err)
        os.Exit(1)
    }
    fmt.Println("Email Sent!")
}
