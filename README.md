# ADservice


<p align="justify">ADservice - маленькая часть сервиса подачи объявлений, созданная в рамках выполнения [тестового задания](https://github.com/avito-tech/verticals/blob/master/trainee/backend.md).</p>

**Объявление** представлено следующей структурой:
```go
type Advert struct {
	ID          int      `json:"id,omitempty"`
	Title       string   `json:"title"`
	Price       int      `json:"price"`
	Date        string   `json:"date,omitempty"`
	Description *string  `json:"description,omitepmty"`
	Gallery     *[]Photo `json:"gallery"`
}
```
_Ограничения_ на поля объявления:
- description (описание) хранит не более 1000 символов;
- title (название) хранит не более 200 символов;
- gallert (галерея фото) хранит не более 3х объектов Photo (фото).

Структура **Photo**, из которых создаётся галерея фото каждого объявления:
```go
type Photo struct {
	Index int    `json:"index"`
	Link  string `json:"photo"`
}
```
Поле index отображает _порядковый номер_ фотографии.
Фото с индексом 0 является **_главным фото_** объявления.


Сервис реализует 3 требуемых метода: 
1. получение списка объявлений;
2. получение конкретного объявления:
3. создание объявления.

... а также необязательные методы:
4. добавление партии обяъвлений;
5. удаление объявления;
