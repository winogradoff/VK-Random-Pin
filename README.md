# go-vk-random-pin

Heroku clock-процесс на Go для закрепление случайного поста на стене ВКонтакте.

### Установка:

```
go get github.com/winogradoff/go-vk-random-pin
```

### Запуск:

```
go-vk-random-pin
```

### В переменных окружения должны быть заданы значения:

```
VK_TOKEN = <токен авторизации>
VK_USERNAME = <имя пользователя>
VK_DELAY = <задержка в секундах>
```

### Они также могут быть заданы в виде ключей при запуске:

```
go-vk-random-pin -token <токен> -username <имя> -delay <задержка>
```

### Токен авторизации можно получить при переходе по следующему URL:

```
https://oauth.vk.com/authorize?client_id=<client_id>&scope=wall,offline&redirect_uri=https://oauth.vk.com/blank.html&display=page&v=5.29&response_type=token
```

где `<client_id>` — идентификатор приложения ВКонтакте.
