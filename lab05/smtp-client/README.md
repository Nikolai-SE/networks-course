
Простой почтовый клиент, позволяющий отправлять письма через TLS tcp-подлючение по сокету.

Нужно использовать пароль для приложений: https://id.yandex.ru/security/app-passwords 
Порт TLS 465 https://yandex.ru/support/mail/mail-clients/ssl.html

```
address   = flag.String("address", "", "email server")
port      = flag.Int("port", 25, "email server port")
sender    = flag.String("sender", "", "sender email")
recipient = flag.String("recipient", "", "recipient email")
text      = flag.String("text", "", "mail text content")
html      = flag.String("html", "", "mail html content")
subject   = flag.String("subject", "", "mail subject")
password  = flag.String("password", "", "password")
```

```bash
go run cli/main/main.go --address smtp.yandex.ru --port 465 --sender "email@mail.com" --recipient "email@mail.com" --subject "Hello" --user user --password password --text "Text"
```

```bash
go run cli/main/main.go --address smtp.yandex.ru --port 465 --sender "email@mail.com" --recipient "email@mail.com" --subject "Hello" --user user --password password --html "<p>Привет <br> <h1> Отправляю html <h1> </p>"
```
