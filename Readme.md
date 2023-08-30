Инструкция по запуску:
Команда 'docker-compose up --build' в консоли в папке проекта

Работа с API:
Адрес запросов:localhost:8000

1. Создание сегмента:
POST Запрос на /segments c данными по типу {"name":"avito-1","percentage":50}
name - имя создаваемого сегмента
percentage - процент пользователей, которые попадут в сегмент автоматически

2. Удаление сегмента:
DELETE Запрос на /segments c данными по типу {"name":"avito-1"}
name - имя удаляемого сегмента

3. Добавление пользователя в сегмент
POST Запрос на /segments/<id пользователя> c данными по типу {"add":["avito1,"avito2"],"remove":["avito3","avito4"]}
add - имена сегментов, которые нужно добавить пользователю
remove - имена сегментов, которые нужно удалить у пользователя
Примечание:
В задании не было указано метода создания пользователя, тем не менее я решил его создать, без предварительного создания пользователя добавить его в сегмент нельзя
POST Запрос на /users с данными по типу {"id":1001,"ttl":2}
id - id создаваемого пользователя
ttl - количество дней, через которое пользователь будет удален

4. Получение сегментов пользователя
GET запрос на /segments с данными по типу {"id":1001}
id - id пользователя, чьи сегменты требуется получить

5. Получение истории удаления и добавления пользователей
GET запрос на /stats с данными по типу {"month":8,"year":2023}
month - месяц, за который нужно вывести историю
year - год, за который нужно вывести историю