# vk-internship
Тестовое задание от VK на позицию Backend разработчик.


# Marketplace Backend API
Тестовое задание на позицию Backend разработчик 

Микросервис для работы с объявлениями. Данный микросервис позволяет создавать, получать, обновлять и удалять задачи, а также позваляет помечать задачу как "выполненной".


Используемые технологии

- golang-migrate/migrate (для миграции)
- pgx (для хранилища данных)
- docker и docker-compose (для запуска сервиса)
- swagger (для API документации)
- gin-gonic/gin (веб фреймворк)

## Запуск

- Заполнить .env как по примеру .env_example
- Запустить сервис с помощью команды `sudo make compose-up`

Документация находится по адресу http://localhost:8000/swagger/index.html с портом 8000 по умолчанию

![alt text](imgs/img-1.png)




## Примеры запросов

- [Регистрация](#sign-up)
- [Вход](#sign-in)
- [Создание объявления](#create-adv)
- [Получение списка объявлений](#get-adv)
- [Получение объявлений пользователя](#get-by-user-id)



### Создание новой задачи <a name="sign-up"></a>

Регистрация

```
curl --location 'http://localhost:8000/auth/sign-up' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "validemail@gmail.com",
  "password": "testtest",
  "name": "test",
  "surname": "test"
}'
```

Пример ответа

```
{
    "id": "3593bccb-5f0d-4c8f-aeca-c42802f68b00"
}
```


### Вход <a name="sign-in"></a>

Пример запроса

```
curl --location 'http://localhost:8000/auth/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "validemail@gmail.com",
  "password": "testtest"
}'
```

Примет ответа

```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE1NTM5NjgsImlhdCI6MTcxMTU0Njc2OCwiVXNlcklEIjoiMzU5M2JjY2ItNWYwZC00YzhmLWFlY2EtYzQyODAyZjY4YjAwIn0.JsV0pL708H8Bgcfx8g2u7ZnrHW-3eF_OFFESpxkCxyo"
}
```

### Создание объявления <a name="create-adv"></a>

Пример запроса 

```
{
  "description": "test",
  "pictures": [
    "https://example.com"
  ],
  "price": 100,
  "title": "test"
}
```

Пример ответа

```
{
    "id": "a14efee4-c271-4e93-a58d-8260f7f42572",
    "code": 201
}
```

### Получение списка объявлений <a name="get-adv"></a>

Пример запроса

```

```

Пример ответа

```

```

### Получение объявлений пользователя <a name="get-by-user-id"></a>

Пример запроса

```


```

Пример ответа

```

```