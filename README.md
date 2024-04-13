# Сервис баннеров
Этот сервис позволяет отображать баннеры для конкретного пользователя на основе
определенных фич и тегов пользователя. Он также предоставляет
функциональные возможности для управления баннерами и связанными с ними тегами
и фичами.

## Обзор:
Сервис работает с баннерами, которые представляют собой JSON-документы,
описывающие элементы пользовательского интерфейса. Каждый баннер связан с
фичей (представленной идентификатором) и несколькими тегами
(представленными идентификаторами). 

### Инструкция по установке
```bash
git clone https://github.com/skraio/banner-service.git
cd banner-service
```

### Создание докер контейнеров
```bash
docker-compose up -d --build
```

### Запуск сервера
```bash
docker-compose up -d
```

### Применение миграций
Скачиваем cli тул [migrate](https://github.com/golang-migrate/migrate/tree/master):
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
Применение миграций
```bash
make migrations-up
```

### Использование API
#### POST /user: Создание пользователя с правами админа:
```bash
curl -d '{"username":"Admin01","password":"strongpassword","role":"admin"}' localhost:8080/v1/user
```
Тело ответа:
```bash
{
    "token": "FSEYMEACLQJTIDJQR5NLUONYAA",
        "user": {
            "user_id": 4,
            "username": "Admin01",
            "role": "admin",
            "created_at": "2024-04-13T13:57:24Z"
        }
}
```

#### POST /banner: Создание баннера с админским токеном:
```bash
BODY='{"tag_ids":[1,2,3],"feature_id":777,"content":{"title":"cakes","text":"homemade birthday cakes","url":"https://example.com"},"is_active":true}'
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" -d "$BODY" localhost:8080/v1/banner
```
Тело ответа:
```bash
{
    "banner_id": 1
}
```

#### POST /user: Создание пользователя с правами пользователя
```bash
curl -i -d '{"username":"User1","password":"strongpassword","role":"user"}' localhost:8080/v1/user
```
Тело ответа:
```bash
{
    "token": "XMIDHIYCER45IU2YTQ726VVH6A",
        "user": {
            "user_id": 3,
            "username": "User01",
            "role": "user",
            "created_at": "2024-04-13T13:55:21Z"
        }
}
```

#### GET /user_banner: Получение баннера по фиче и тегу
```bash
curl -H "Authorization: Bearer XMIDHIYCER45IU2YTQ726VVH6A" 'localhost:8080/v1/user_banner?feature_id=777&tag_id=2'
```
Тело ответа:
```bash
{
    "content": {
        "title": "cakes",
            "text": "homemade birthday cakes",
            "url": "https://example.com"
    }
}
```

#### GET /banner: Получение списка баннеров по фиче и/или тегу
1) По фиче:
```bash
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" 'localhost:8080/v1/banner?feature_id=777'
```
Тело ответа:
```bash
{
    "banners": [
    {
        "banner_id": 1,
            "feature_id": 777,
            "tag_ids": [
                1,
                2,
                3
            ],
            "content": {
                "title": "cakes",
                "text": "homemade birthday cakes",
                "url": "https://example.com"
            },
            "is_active": true,
            "created_at": "2024-04-13T14:00:37.94Z",
            "updated_at": "2024-04-13T14:00:37.94Z"
    },
    {
        "banner_id": 2,
        "feature_id": 777,
        "tag_ids": [
            4,
            5,
            6
        ],
        "content": {
            "title": "luxury sofas",
            "text": "high quality sofas",
            "url": "https://example.com"
        },
        "is_active": true,
        "created_at": "2024-04-13T14:09:10.664Z",
        "updated_at": "2024-04-13T14:09:10.664Z"
    }
    ],
    "metadata": {
        "Offset": 0,
        "Limit": 10,
        "TotalRecords": 2
    }
}
```

2) По тегу:
```bash
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" 'localhost:8080/v1/banner?tag_id=1'
```
Тело ответа:
```bash
{
    "banners": [
    {
        "banner_id": 1,
            "feature_id": 777,
            "tag_ids": [
                1,
                2,
                3
            ],
            "content": {
                "title": "cakes",
                "text": "homemade birthday cakes",
                "url": "https://example.com"
            },
            "is_active": true,
            "created_at": "2024-04-13T14:00:37.94Z",
            "updated_at": "2024-04-13T14:00:37.94Z"
    },
    {
        "banner_id": 3,
        "feature_id": 999,
        "tag_ids": [
            1,
            5,
            9
        ],
        "content": {
            "title": "kitchen tables",
            "text": "enjoy your meals",
            "url": "https://example.com"
        },
        "is_active": true,
        "created_at": "2024-04-13T14:09:18.032Z",
        "updated_at": "2024-04-13T14:09:18.032Z"
    }
    ],
    "metadata": {
        "Offset": 0,
        "Limit": 10,
        "TotalRecords": 2
    }
}
```

3) По фиче и тегу:
```bash
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" 'localhost:8080/v1/banner?tag_id=1&feature_id=999'
```
Тело ответа:
```bash
{
    "banners": [
    {
        "banner_id": 3,
            "feature_id": 999,
            "tag_ids": [
                1,
                5,
                9
                ],
            "content": {
                "title": "kitchen tables",
                "text": "enjoy your meals",
                "url": "https://example.com"
            },
            "is_active": true,
            "created_at": "2024-04-13T14:09:18.032Z",
            "updated_at": "2024-04-13T14:09:18.032Z"
    }
    ],
    "metadata": {
        "Offset": 0,
        "Limit": 10,
        "TotalRecords": 1
    }
}
```

#### PATCH /banner/{id}: Обновление содержимого баннера
```bash
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" -X PATCH -d '{"feature_id":90909,"is_active":false,"content":{"text":"stylish and functional kitchen tables"}}' "localhost:8080/v1/banner/3"
```
Тело ответа:
```bash
{
    "banner": {
        "banner_id": 3,
            "feature_id": 90909,
            "tag_ids": [
                1,
                5,
                9
            ],
            "content": {
                "title": "kitchen tables",
                "text": "stylish and functional kitchen tables",
                "url": "https://example.com"
            },
            "is_active": false,
            "created_at": "2024-04-13T14:09:18.032Z",
            "updated_at": "2024-04-13T14:27:24.354Z"
    }
}
```

#### DELETE /banner/{id}: Удаление баннера
```bash
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" -X DELETE "localhost:8080/v1/banner/2"
curl -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" "localhost:8080/v1/banner"
```
Результат - баннер успешно удален
```bash
{
        "banners": [
                {
                        "banner_id": 1,
                        "feature_id": 777,
                        "tag_ids": [
                                1,
                                2,
                                3
                        ],
                        "content": {
                                "title": "cakes",
                                "text": "homemade birthday cakes",
                                "url": "https://example.com"
                        },
                        "is_active": true,
                        "created_at": "2024-04-13T14:00:37.94Z",
                        "updated_at": "2024-04-13T14:00:37.94Z"
                },
                {
                        "banner_id": 3,
                        "feature_id": 90909,
                        "tag_ids": [
                                1,
                                5,
                                9
                        ],
                        "content": {
                                "title": "kitchen tables",
                                "text": "stylish and functional kitchen tables",
                                "url": "https://example.com"
                        },
                        "is_active": false,
                        "created_at": "2024-04-13T14:09:18.032Z",
                        "updated_at": "2024-04-13T14:27:24.354Z"
                }
        ],
        "metadata": {
                "Offset": 0,
                "Limit": 10,
                "TotalRecords": 2
        }
}
```
### Нагрузочное тестирование
При 1000 запросах среднее время ответа 4ms и SLI успешности ответа 1000/1000
```bash
hey -m GET \
    -H "Authorization: Bearer FSEYMEACLQJTIDJQR5NLUONYAA" \
    -n 1000 \
    "http://localhost:8080/v1/user_banner?tag_id=1&feature_id=90909"


Summary:
  Total:        0.0866 secs
  Slowest:      0.0127 secs
  Fastest:      0.0003 secs
  Average:      0.0040 secs
  Requests/sec: 11541.7723

  Total data:   132000 bytes
  Size/request: 132 bytes

Response time histogram:
  0.000 [1]     |
  0.002 [70]    |■■■■■■■■
  0.003 [209]   |■■■■■■■■■■■■■■■■■■■■■■■
  0.004 [360]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [154]   |■■■■■■■■■■■■■■■■■
  0.007 [79]    |■■■■■■■■■
  0.008 [56]    |■■■■■■
  0.009 [32]    |■■■■
  0.010 [28]    |■■■
  0.011 [7]     |■
  0.013 [4]     |


Latency distribution:
  10% in 0.0018 secs
  25% in 0.0026 secs
  50% in 0.0034 secs
  75% in 0.0049 secs
  90% in 0.0069 secs
  95% in 0.0086 secs
  99% in 0.0104 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0003 secs, 0.0127 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0035 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0027 secs
  resp wait:    0.0038 secs, 0.0003 secs, 0.0118 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0042 secs

Status code distribution:
  [200] 1000 responses
```
