# Go Package Manager CLI (`pm`)

`pm` — кросс-платформенный CLI-инструмент на Go 1.24+ для упаковки файлов в архивы и их распространения через SSH/SFTP.

## Возможности
* **pm create** — собирает файлы по конфигурации (`packet.json|yaml`), формирует ZIP `name-ver.zip` и при необходимости загружает на сервер.
* **pm update** — скачивает архивы по конфигурации (`packages.json|yaml`), проверяет версии (semver), распаковывает в указанную директорию.

## Установка
```bash
# клонировать
 git clone https://github.com/your-org/package-manager.git
 cd package-manager

# собрать
 go install ./cmd/pm       # бинарник в $GOPATH/bin
```

## Форматы конфигураций
### packet.json
```json
{
  "name":   "packet-1",
  "ver":    "1.10",
  "targets": [
    "./src/**/*.go",
    { "path": "./assets", "exclude": "*.tmp" }
  ]
}
```

### packages.json
```json
{
  "packages": [
    { "name": "packet-1", "ver": ">=1.10" }
  ]
}
```

## Быстрый старт
1. Сгенерируйте SSH-ключ **без пароля** и скопируйте `id_rsa.pub` на сервер (`~/.ssh/authorized_keys`).
2. Создайте `packet.json`, затем:
   ```powershell
   pm create packet.json \
     --host myserver:22 --user dev \
     --key "$Env:USERPROFILE\.ssh\id_rsa" \
     --remote-dir /opt/pm_repo
   ```
3. На целевой машине:
   ```powershell
   pm update packages.json \
     --host myserver:22 --user dev \
     --key "$Env:USERPROFILE\.ssh\id_rsa" \
     --remote-dir /opt/pm_repo \
     --dest .\vendor
   ```

## Сборка и отладка в VS Code / Windsurf IDE
В репозитории есть `.vscode/launch.json` с готовыми конфигурациями "pm create" и "pm update".

## Разработка
* Код разбит на модули: `cmd`, `internal/config`, `matcher`, `archiver`, `transport`, `version`.
* Зависимости управляются через Go modules (см. `go.mod`).
* Для тестов используйте `go test ./...`.
