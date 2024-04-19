
Простой почтовый клиент, позволяющий отправлять письма через консоль.

Нужно использовать пароль для приложений:
https://id.yandex.ru/security/app-passwords 


```
address   = flag.String("address", "", "email server")
port      = flag.Int("port", 25, "email server port")
sender    = flag.String("sender", "", "sender email")
recipient = flag.String("recipient", "", "recipient email")
text      = flag.String("text", "", "mail text content")
html      = flag.String("html", "", "mail html content")
subject   = flag.String("subject", "", "mail subject")
user      = flag.String("user", "", "user")
password  = flag.String("password", "", "password")
```

```bash
go run main/main.go --address smtp.yandex.ru --port 587 --sender "email@mail.com" --recipient "email@mail.com" --subject "Hello" --user user --password password --text "Text"
```

```bash
go run main/main.go --address smtp.yandex.ru --port 587 --sender "email@mail.com" --recipient "email@mail.com" --subject "Hello" --user user --password password --html "<p>Привет <br> <h1> Отправляю html <h1> </p>"
```
