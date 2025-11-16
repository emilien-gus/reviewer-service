# reviewer-service

## Описание

Сервис для работы с пулл реквестами и ревьюерами.  

## Установка и запуск

Для запуска проекта вам потребуется [Docker](https://www.docker.com/get-started) и [Docker Compose](https://docs.docker.com/compose/install/).

### 1. Клонирование репозитория

```bash
git clone https://github.com/emilien-gus/reviewer-service
cd reviewer-service
```

### 2. Запуск проекта
Создайте .env файл и заполните его переменными из .env.example.
Пример содержания .env файла:
```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=pr_reviewer
```
Далее выполните следующую команду
```bash
docker-compose up --build
```

### 3. Дополнительные задания
Добавил endpoint статистики (GET /stats), которую возвращает количество назначений по пользователям и по PR.

### 4. Вопросы и решения
1.   тэг - name: Health упоминается в api, но нигде не используется. Зачем он? Решил проигнорировать его, так как нет описания, что с ним связано. 
