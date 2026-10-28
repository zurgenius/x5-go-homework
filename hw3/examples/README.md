# Примеры занятия 5

Запускайте из папки `student`, рядом с `go.mod`. Каждый пример независим от незаконченной практики.

| Папка | Что показывает | Команда |
|---|---|---|
| `01-lesson4` | Сервис с прошлого занятия | `go run ./examples/01-lesson4` |
| `02-error-value` | Два результата функции | `go run ./examples/02-error-value` |
| `03-sentinel-wrap` | Sentinel, wrapping, Is и Unwrap | `go run ./examples/03-sentinel-wrap` |
| `04-custom-error` | Поля собственного типа ошибки | `go run ./examples/04-custom-error` |
| `05-defer-order` | Порядок вызовов defer | `go run ./examples/05-defer-order` |
| `06-defer-resource` | Закрытие ресурса при раннем return | `go run ./examples/06-defer-resource` |
| `07-panic` | Паника без recover | `go run ./examples/07-panic` |
| `08-recover-boundary` | Именованный результат и recover | `go run ./examples/08-recover-boundary` |
| `09-recover-loop` | Почему recover вокруг цикла недостаточно | `go run ./examples/09-recover-loop` |
| `10-student-example` | Готовый пример для самостоятельного повторения | `go run ./examples/10-student-example` |

Пример `07-panic` намеренно завершается с паникой и ненулевым кодом. Это демонстрация, а не ошибка установки. Следующий пример показывает перехват.

Для самостоятельного повторения откройте `10-student-example`. Запустите его, затем добавьте ещё один товар и ещё один запрос. Проверьте, что количество 0 и отсутствие товара обрабатываются по-разному. Подробности есть в README этой папки.
