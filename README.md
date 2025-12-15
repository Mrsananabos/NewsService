# News REST API

REST API сервис для управления новостями на Go + Fiber + PostgreSQL + Reform.

## Технологии

- **Go 1.25** 
- **Fiber** - веб-фреймворк
- **Reform** - ORM для работы с БД
- **PostgreSQL** - база данных
- **Goose** - миграции
- **Logrus** - логирование
- **Docker** - контейнеризация
- **Swagger** - документация API

### 1. Клонировать репозиторий
```bash
git clone <repository-url>
cd <project-directory>
```

### 2. Настроить переменные окружения
```env
PORT=8080
DB_ADDRESS=postgresql
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
BEARER_TOKEN=my-secret-token-9999  
```

### 3. Запустить через Docker Compose
```bash
docker-compose up --build
```

Сервис будет доступен по адресу: `http://localhost:8080`

## API Endpoints

### Аутентификация
Все запросы требуют Bearer токен в заголовке:
```
Authorization: Bearer my-secret-token-9999
```
### 1. Создание новости
```http
POST /create
Content-Type: application/json

{
  "Title": "New Title",
  "Content": "New Content",
  "Categories": [1, 2]
}
```

**Обязательные поля:**
- `Title` (макс. 255 символов)
- `Content`

**Ответ:**
```json
{
  "Success": true,
  "Id": 1
}
```


### 2. Редактирование новости
```http
POST /edit/:id
Content-Type: application/json

{
  "Title": "Updated Title",
  "Content": "Updated Content",
  "Categories": [1, 2, 3]
}
```

**Особенности:**
- Все поля опциональны
- Обновляются только переданные поля
- `Categories` полностью заменяет существующие

**Ответы:**
- `200` - успешно обновлено
- `400` - ошибка валидации
- `401` - неверный токен
- `404` - новость не найдена

### 3. Список новостей
```http
GET /list?limit=10&offset=0
```

**Параметры:**
- `limit` (опционально) - количество записей (1-100, по умолчанию 10)
- `offset` (опционально) - смещение (по умолчанию 0)

**Ответ:**
```json
{
  "success": true,
  "News": [
    {
      "Id": 1,
      "Title": "News Title",
      "Content": "News Content",
      "Categories": [1, 2, 3]
    }
  ]
}
```

## Документация API (Swagger)

После запуска сервиса откройте:
```
http://localhost:8080/swagger/index.html
```
## Структура БД

### Таблица `news`
```sql
id       BIGSERIAL PRIMARY KEY
title    VARCHAR(255) NOT NULL
content  TEXT NOT NULL
```

### Таблица `news_categories`
```sql
news_id      BIGINT NOT NULL
category_id  BIGINT NOT NULL
PRIMARY KEY (news_id, category_id)
FOREIGN KEY (news_id) REFERENCES news(id)
```