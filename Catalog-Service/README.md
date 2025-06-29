# Catalog Service

## API

### REST:
- `POST /items` - создание предмета в каталоге.
---
- `GET /items` - получение предметов в каталоге.

|  Query   |       Value      |
|----------|------------------|
| category | `{category_name}` |
---
- `GET /items/{item_id}` - получение предмета в каталоге
---
- `DELETE /items/{item_id}` - удаление товара из каталога.
---

### gRPC(CatalogService):
- `GetItem(ID) Item` - получение товара (`Item`) по идентификатору (`ID`).