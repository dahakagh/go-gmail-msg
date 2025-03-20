# go-gmail-msg

# 1️⃣ Установка зависимостей

Убедитесь, что у вас установлен Go 1.18+. Затем клонируйте репозиторий и установите зависимости:

```sh
git clone https://github.com/yourusername/gmail-reader.git
cd gmail-reader
go mod tidy
```

# 2️⃣ Настройка Google API

1. Перейдите в Google Cloud Console и создайте новый проект.
2. Включите Gmail API в разделе APIs & Services > Library.
3. Создайте OAuth 2.0 Client ID в APIs & Services > Credentials.
4. В разделе Authorized redirect URIs укажите:

   ```
    http://localhost:8080/callback
   ```

5. Скачайте credentials.json и сохраните его в переменные окружения (см. ниже).

# 3️⃣ Настройка .env

Создайте файл .env в корне проекта:

```
GMAIL_CLIENT_ID=your-client-id
GMAIL_CLIENT_SECRET=your-client-secret
GMAIL_REDIRECT_URI=your-redirect-uri
HTTP_SERVER_URL=localhost:8080
```

Важно! Добавьте .env в .gitignore, чтобы не публиковать API-ключи.

# 4️⃣ Запуск сервиса

```sh
go run main.go
```

При первом запуске сервис выдаст ссылку для авторизации. Перейдите по ней, войдите в Google.
Токен OAuth будет автоматически сохранен в token.json, и повторная авторизация не потребуется.

# 📁 Где хранятся письма и вложения?

Письма: emails/
Вложения: attachments/
Токен Google: token.json
