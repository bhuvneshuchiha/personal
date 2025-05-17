package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type Notifier interface {
	Notify() string
}

type EmailUser struct {
	Name    string
	Contact string
}

type SMSUser struct {
	Name    string
	Contact string
}

func (e EmailUser) Notify() string {
	templ, er := template.New("test").Parse("Email sent to {{.Name}} at {{.Contact}}")
	if er != nil {
		return er.Error()
	}
	var buf bytes.Buffer
	templ.Execute(&buf, e)
	return buf.String()
}

func (s SMSUser) Notify() string {
	templ, er := template.New("my").Parse("SMS sent to {{.Name}} at {{.Contact}}")
	if er != nil {
		return er.Error()
	}
	var buf bytes.Buffer
	templ.Execute(&buf, s)
	return buf.String()
}

func SendNotification(n Notifier) {
	notString := n.Notify()
	fmt.Println("The output from interface implementor:-", notString)
}

func main() {
	e := EmailUser{"Alice", "alice@example.com"}
	s := SMSUser{"Bob", "+911234567890"}

	SendNotification(e)
	SendNotification(s)
}
