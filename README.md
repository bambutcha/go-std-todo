# Todo HTTP Server

HTTP сервер на Go для управления задачами (todos) с CRUD операциями.

## Требования

- Go 1.25.3 или выше
- Docker (опционально, для запуска в контейнере)

## Запуск приложения

### Локальный запуск

```bash
go run ./cmd
```

Или скомпилировать и запустить:

```bash
go build -o todo-server ./cmd
./todo-server
```

Сервер запустится на порту `8080` по умолчанию. Можно изменить порт через переменную окружения:

```bash
PORT=3000 ./todo-server
```

Или создать файл `.env` в корне проекта:

```bash
PORT=3000
```

Приложение автоматически загрузит переменные из `.env` файла (если он существует).

### Запуск через Docker

```bash
docker build -t todo-server .
docker run -p 8080:8080 todo-server
```

## API Endpoints

- `POST /todos` - создать новую задачу
- `GET /todos` - получить список всех задач
- `GET /todos/{id}` - получить задачу по идентификатору
- `PUT /todos/{id}` - обновить задачу по идентификатору
- `DELETE /todos/{id}` - удалить задачу по идентификатору

## Структура задачи

```json
{
  "id": 1,
  "title": "Task title",
  "description": "Task description",
  "completed": false
}
```

## Примеры использования

### Создать задачу

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "My first task", "description": "Task description", "completed": false}'
```

### Получить все задачи

```bash
curl http://localhost:8080/todos
```

### Получить задачу по ID

```bash
curl http://localhost:8080/todos/1
```

### Обновить задачу

```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated title", "description": "Updated description", "completed": true}'
```

### Удалить задачу

```bash
curl -X DELETE http://localhost:8080/todos/1
```

## Тестирование

Запустить все тесты:

```bash
go test ./...
```

## Особенности

- In-memory хранилище данных (без внешних БД)
- Валидация: заголовок не может быть пустым
- Обработка ошибок: 400 Bad Request для валидации, 404 Not Found для несуществующих задач
- Логирование всех HTTP запросов
- Использование context для таймаутов (5 секунд)
- Graceful shutdown при получении SIGINT/SIGTERM
- Только стандартная библиотека Go (без сторонних зависимостей)
