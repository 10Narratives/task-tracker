# Task Tracker

## Description

Task tracker is a simple web application which provides an opportunity to manage your TODO tasks

## Installation

Install sqlite on your computer

```shell
sudo pacman -S sqlite
```

Clone this repository on your computer

```shell
git clone git@github.com:10Narratives/task-tracker.git
```

Create `.env` file

```.env
CONFIG_PATH=/path/to/your/config
```

Write configuration file in `yaml` format

```yaml
env: "local" # logger option
storage: # database settings
  driver: "sqlite3"
  dsn: "storage/scheduler.db"
  limit: 10
http_server: # http server settings
  address: "localhost"
  port: "8000"
  timeout: 4s
  idle_timeout: 60s
  file_server_path: "./web"
```

## Usage

## License
