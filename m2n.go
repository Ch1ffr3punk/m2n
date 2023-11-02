package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "net/smtp"
    "os"
)

func main() {
    // contacting server as
    from := "user@localhost"

    // sending message to
    to := "mail2news@dizum.com"

    // server we send message through
    host := "127.0.0.1"

    // check if there is piped input
    fi, _ := os.Stdin.Stat()
    if (fi.Mode() & os.ModeCharDevice) != 0 {
        fmt.Println("Usage: m2n < message.txt") // with  all Usenet message Headers
        os.Exit(1)
    }

    message, _ := ioutil.ReadAll(os.Stdin)

    // Connect to the server, set the sender and recipient,
    // and send the message all in one step.
    
    c, err := smtp.Dial(host + ":5870")
    if err != nil {
        fmt.Println("Error Dial: ", err)
        os.Exit(1)
    }
    
    c.StartTLS(&tls.Config{InsecureSkipVerify: true})
    
    if err = c.Mail(from); err != nil {
        fmt.Println("Error Mail: ", err)
        os.Exit(1)
    }
    
    if err = c.Rcpt(to); err != nil {
        fmt.Println("Error Rcpt: ", err)
        os.Exit(1)
    }
    
    w, err := c.Data()
    if err != nil {
        fmt.Println("Error Data: ", err)
        os.Exit(1)
    }
    
    _, err = w.Write(message)
    if err != nil {
        fmt.Println("Error Write: ", err)
        os.Exit(1)
    }
    
    err = w.Close()
    if err != nil {
        fmt.Println("Error Close: ", err)
        os.Exit(1)
    }
    
    c.Quit()

    fmt.Println("Message Sent!")
}
