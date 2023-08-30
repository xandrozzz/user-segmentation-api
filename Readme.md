Инструкция по запуску:
Команда 'docker-compose up --build' в консоли в папке проекта

Работа с API:
Адрес запросов: localhost:8000

1. Создание сегмента:
POST Запрос на /segments c данными по типу {"name":"avito-1","percentage":50}
name - имя создаваемого сегмента
percentage - процент пользователей, которые попадут в сегмент автоматически

Пример:

![Screenshot_61](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/1e39387f-b6b7-4634-9192-b6ac3f2fa3e3)

2. Удаление сегмента:
DELETE Запрос на /segments c данными по типу {"name":"avito-1"}
name - имя удаляемого сегмента

Пример:

![Screenshot_62](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/94c0cee7-2d12-4eea-8654-c8156ba3291a)

3. Добавление пользователя в сегмент
POST Запрос на /users/<id пользователя> c данными по типу {"add":["avito1,"avito2"],"remove":["avito3","avito4"]}
add - имена сегментов, которые нужно добавить пользователю
remove - имена сегментов, которые нужно удалить у пользователя

Пример:

![Screenshot_64](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/70e0a1c0-e280-4c55-972d-2308c579fc61)

Примечание:
В задании не было указано метода создания пользователя, тем не менее я решил его создать, без предварительного создания пользователя добавить его в сегмент нельзя
POST Запрос на /users с данными по типу {"id":1001,"ttl":2}
id - id создаваемого пользователя
ttl - количество дней, через которое пользователь будет удален

Пример:

![Screenshot_58](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/99ee0809-f56e-43c5-8c53-671a869556ad)

4. Получение сегментов пользователя
GET запрос на /segments с данными по типу {"id":1001}
id - id пользователя, чьи сегменты требуется получить

Пример:

![Screenshot_63](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/d390a274-447a-4d36-92d8-8d8468b7bbf7)

5. Получение истории удаления и добавления пользователей
GET запрос на /stats с данными по типу {"month":8,"year":2023}
month - месяц, за который нужно вывести историю
year - год, за который нужно вывести историю

Пример:

![Screenshot_65](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/9dadab9c-6b8e-4da8-a8de-328122bc2e19)

6. Дополнительный метод удаления пользователей
DELETE Запрос на /users c данными по типу {"id":1001}
id - имя удаляемого сегмента

Пример:

![Screenshot_66](https://github.com/xandrozzz/user-segmentation-api/assets/94150105/6817dd0f-db82-4e6f-b4d3-353911c5f2f7)
