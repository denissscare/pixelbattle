# PixelBattle Backend

Реализация серверной части для онлайн-игры PixelBattle — с регистрацией, авторизацией, WebSocket-канвасом, системой пикселей, мониторингом и полной контейнеризацией.

## Описание проекта

PixelBattle — учебный проект, реализующий современный backend на языке Go. В качестве базы данных используется PostgreSQL, для хранения файлов — Minio (аналог S3), для кэширования и быстрой работы — Redis, а также для real-time обмена событиями используется NATS. Весь проект упакован в Docker-контейнеры и сопровождается мониторингом через Prometheus и Grafana.

## Технологии

- **Go** — основной язык приложения.
- **chi** — фреймворк для REST и WebSocket API.
- **PostgreSQL** — реляционная база данных для хранения пользователей и истории пикселей .
- **Redis** — хранилище для актуального canvas пикселей.
- **Minio** — файловое хранилище (S3) для аватаров.
- **NATS** — брокер сообщений для real-time взаимодействия.
- **Prometheus & Grafana** — сбор и визуализация метрик.
- **Docker Compose** — контейнеризация и быстрое развертывание.

## Функциональность

| Метод | Маршрут         | Описание                                                  | Аутентификация |
| ----- | --------------- | --------------------------------------------------------- | -------------- |
| POST  | /register       | Регистрация пользователя + загрузка аватара (опционально) | Нет            |
| POST  | /login          | Авторизация, получение JWT (кука)                         | Нет            |
| POST  | /avatar         | Загрузка/обновление аватара пользователя                  | Да (JWT)       |
| POST  | /email          | Изменение почты пользователя                              | Да (JWT)       |
| GET   | /canvas         | Получение текущего состояния канваса                      | Да (JWT)       |
| POST  | /pixel          | Отправка нового пикселя (участие в битве)                 | Да (JWT)       |
| GET   | /pixels/history | Получение истории изменений пикселей                      | Да (JWT)       |
| GET   | /ws             | WebSocket-подключение к канвасу (real-time updates)       | Нет            |
| GET   | /metrics        | Prometheus-метрики                                        | Нет            |

## Требования

- **Go** 1.22+
- **Docker** и **Docker Compose**
- **Git** (для клонирования)

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone git@gitlab.com:pixelbattle1/pixelbattle.git
cd pixelbattle
```

### 2. Подготовка переменных окружения

Создайте файл `.env` в корне проекта со следующим содержимым:

```env
REDIS_PASSWORD=your_redis_password
REDIS_USER=your_redis_user
REDIS_USER_PASSWORD=your_redis_user_password
GRAFANA_PASSWORD=your_grafana_password

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=your_pg_user
POSTGRES_PASSWORD=your_pg_password
POSTGRES_DB=pixelbattle
POSTGRES_SSLMODE=disable

MINIO_ENDPOINT=minio:9000
MINIO_PUBLIC_HOST=localhost:9000
MINIO_ACCESS_KEY=your_minio_access_key
MINIO_SECRET_KEY=your_minio_secret_key
MINIO_BUCKET=avatars
MINIO_USE_SSL=false
```

### 3. Запуск контейнеров

```bash
docker compose up --build -d
```

- Все сервисы будут развернуты автоматически: **Postgres**, **Redis**, **NATS**, **Minio**, **Prometheus**, **Grafana**, и Go-сервер.
- После запуска:
  - API доступен на [http://localhost:8080](http://localhost:8080)
  - Grafana — [http://localhost:3001](http://localhost:3001) (логин: admin / пароль из переменной `GRAFANA_PASSWORD`)
  - Prometheus — [http://localhost:9090](http://localhost:9090)
  - Minio — [http://localhost:9001](http://localhost:9001) (логин/пароль из переменных)

### 4. Посмотреть логи сервера

```bash
docker compose logs server
```

или

```bash
docker compose logs -f server
```

чтобы показывать логи в реальном времени

### 5. Остановка

```bash
docker compose down
```

или

```bash
docker compose down -v
```

чтобы удалить все volumes

## Примеры API-запросов

### Регистрация

```bash
curl -X POST -F "username=testuser" -F "email=test@e.com" -F "password=1234" \
     -F "avatar=@/path/to/avatar.jpg" http://localhost:8080/register
```

(аватар опционально, можно не указывать)

### Авторизация

```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"email":"test@e.com","password":"1234"}' \
     http://localhost:8080/login
```

ответ:

```json
{
  "token": "jwt-токен",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@e.com",
    "avatar": "http://localhost:9000/avatars/1.jpg"
  }
}
```

### WebSocket

Получение изменения пикселей в реальном времени + получение последнего пикселя пользователя + получение актуального канваса при подключении

```bash
wscat -c "ws://0.0.0.0:8080/ws?username=testuser"
```

## Тестирование

Проект включает unit-тесты для проверки сервисов auth и pixelbattle:

1. Перейдите в корневую директори:
2. Запуск тестов:

```shell
go test ./...
```

Тесты проверяют:

### **auth (user service)**

- **TestRegisterUserAlreadyExists:** Проверяет ошибку при регистрации пользователя с уже существующим email.
- **TestRegisterUsernameExists:** Проверяет ошибку при регистрации пользователя с уже существующим username.
- **TestRegisterSuccess:** Проверяет успешную регистрацию нового пользователя.
- **TestLoginByEmailSuccess:** Проверяет успешный вход по email и валидному паролю.
- **TestLoginInvalidPassword:** Проверяет ошибку при попытке входа с неверным паролем.
- **TestRegisterWithIDSuccess:** Проверяет успешную регистрацию с возвратом ID нового пользователя.
- **TestUpdateAvatarURLSuccess:** Проверяет успешное обновление ссылки на аватар пользователя.
- **TestUploadAvatarSuccess:** Проверяет успешную загрузку аватара.
- **TestUploadAvatarS3Error:** Проверяет обработку ошибки при загрузке аватара в S3.
- **TestUploadAvatarDBError:** Проверяет обработку ошибки при обновлении ссылки аватара в БД.
- **TestUpdateEmailSuccess:** Проверяет успешное обновление email пользователя.
- **TestUpdateEmailEmailExists:** Проверяет ошибку при попытке изменить email на уже занятый.
- **TestUpdateEmailRepoError:** Проверяет обработку ошибки при обновлении email пользователя в БД.

### **config**

- **TestLoadConfig:** Проверяет корректную загрузку конфигурации приложения из yaml-файла.

### **pixcelbattle/service (battle service)**

- **TestInitCanvasM:** Проверяет получение канваса из Redis и обработку ошибки.
- **TestUpdatePixelSuccess:** Проверяет успешное обновление пикселя — сохранение в Redis, публикация через брокер, инкремент метрики.
- **TestUpdatePixelValidationError:** Проверяет обработку ошибки валидации при попытке установить некорректный пиксель.
- **TestUpdatePixelRedisError:** Проверяет обработку ошибки записи пикселя в Redis.
- **TestUpdatePixelBrokerError:** Проверяет обработку ошибки публикации обновления пикселя через брокер (NATS).

### **hash**

- **TestHashPassword:** Проверяет корректность хеширования пароля и его успешную проверку через bcrypt.

### **domain**

- **TestValidate:**  
   Проверяет корректную работу метода `Validate()` для структуры `Pixel`:
  - Ошибка при X вне допустимого диапазона.
  - Ошибка при пустом поле автора.
  - Ошибка при некорректной дате (в будущем).
