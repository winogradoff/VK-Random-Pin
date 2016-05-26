# go-vk-random-pin

Heroku приложение на Go для закрепление случайного поста на стене ВКонтакте.

### Установка:

```
go get -u github.com/winogradoff/go_vk_random_pin/...
```

### Запуск:

```
go_vk_random_pin_web
go_vk_random_pin_worker
```

### В переменных окружения должны быть заданы значения:

```
VK_TOKEN = <токен авторизации>
VK_USERNAME = <имя пользователя>
VK_DELAY = <задержка в секундах>
DATABASE_URL = <строка подключения к БД>
PORT = <порт (автоматически создается Heroku)>
```

### Токен авторизации можно получить при переходе по следующему URL:

```
https://oauth.vk.com/authorize?client_id=<client_id>&scope=wall,offline&redirect_uri=https://oauth.vk.com/blank.html&display=page&v=5.29&response_type=token
```

где `<client_id>` — идентификатор приложения ВКонтакте.
