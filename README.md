# Product Catalog
Сервис для управления товарами

## Запуск

#### Docker контейнеры
```
make up service=auth

make up service=app

```

#### Запуск сервисов
```
go run .\cmd\app\main.go

go run .\cmd\auth\main.go
```

#### Запуск Swagger
Через веб браузер
В 
```http://localhost:8080/swagger/index.html```
Логинимся получаем Access Token

Затем переходим в
```http://localhost:8000/swagger/index.html```
проходим аутефикацию (Bearer <Token>)


## API Эндпоинты

```http request
GET    api/v1/health/            // проверка работоспособности сервиса
POST   api/v1/auth/register      // регистрация
POST   api/v1/auth/login         // логин
POST   api/v1/auth/refresh       // обновление JWT токенов

POST    api/v1/categories         // добавить категорию 
GET     api/v1/categories         // Получить все категории
GET     api/v1/categories/stats   // Получить статистику
PUT     api/v1/categories/:id     // Изменить категорию
DELETE  api/v1/categories/:id     // Удалить категорию

POST    api/v1/products                 // добавить товар
GET     api/v1/products                 // Получить все товары
GET     api/v1/products/:id             // Получить товар по ID
PUT     api/v1/products/:id             // Изменить товар
DELETE  api/v1/products/:id             // Удалить товар
PUT     api/v1/products/:id/restore     // Восстановить товар
PUT     api/v1/products/:id/status      // Изменить статус товара
POST    api/v1/products/:id/image       // Загрузить изображение товара
GET     api/v1/products/:id/image       // Получить изображение товара
```
