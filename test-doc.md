# Временная документация

Каждый ответ будет содержать в себе JSON строку:

```json
{"success": true, "response": [...]}
```

При ошибках будет выдан такой ответ:

```json
{"success": false, "erron": "..."}
```

## Авторизация

### Регистрация

`POST /api/register`
Обязательные параметры:
- `first_name` - Имя
- `last_name` - Фамилия
- `middle_name` - Отчество
- `email` - E-mail
- `password` - Пароль

### Авторизация

`POST /api/login`
Обязательные параметры:
- `email` - E-mail
- `password` - Пароль

### Ответ при успехе:

```json
{"success": true, "response": {"email": "<EMAIL_ПОЛЬЗОВАТЕЛЯ>", "access_token": "<ТУТ_ТОКЕН>"}}
```

## Создание заявки
