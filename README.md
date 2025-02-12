# Task Tracker

## 📝 Description

Task tracker is a simple web application which provides an opportunity to manage your TODO tasks

## ✨ Features

### 🔄  **Automatic Task Rescheduling**  

The application supports a flexible rule system for automatically rescheduling completed tasks.  

#### 🗑️ **Deleting Completed Tasks**  

If no rescheduling rule is specified, the completed task will be **removed** from the list.  

#### 📅 **Task Rescheduling Rules**  

| Rule                    | Description                                                                       |
| ----------------------- | --------------------------------------------------------------------------------- |
| `d <number>`            | Moves the task forward by the specified number of days (max. 400)                 |
| `y`                     | Reschedules the task for the **same date next year**                              |
| `w <1-7>`               | Assigns the task to the nearest specified weekday *(1 — Mon, 7 — Sun)*            |
| `m <1-31,-1,-2> [1-12]` | Assigns the task to specific days of the month, optionally within specific months |

### 🔎 Search and Filtering  

The application provides two ways to find tasks:  

#### 🔍 **Search by Content**  

Users can enter a **keyword** in the search field next to the "Add Task" button.  
The system will look for this keyword in the **title, comments, and content** of tasks.  

#### 📅 **Filter by Date**  

Users can enter a **specific date** in the format `DD.MM.YYYY` to filter tasks.  
The system will return only those tasks that are scheduled for the given date.  

### 🔐 **Authentication with JWT**

The application uses JSON Web Tokens (JWT) for secure authentication.

## 📌 Prerequisites  

Before setting up the project, ensure you have the following installed:  

- **Go (Golang)**: [Download and install Go](https://golang.org/dl/)  
- **SQLite**: [Download and install SQLite](https://www.sqlite.org/download.html)  

## 🚀 Installation  

Follow these steps to set up the project:  

### 1️⃣ **Clone the Repository**  

```bash
git clone https://github.com/10Narratives/task-tracker.git
cd task-tracker
```

### 2️⃣ Install Dependencies

Run the following command to install required Go packages:

```bash
go get ./...
```

### 3️⃣ Set Up Environment Variables

Create a .env file and specify the configuration path:

```dotenv
CONFIG_PATH=./path/to/your/config.yaml
PASSWORD=your_password
```

### 4️⃣ Create a Custom Configuration File

The project uses a YAML configuration file to manage environment and service settings. Below is a breakdown of available configuration parameters:

| Parameter                      | Type   | Description                                   | Default value            |
| ------------------------------ | ------ | --------------------------------------------- | ------------------------ |
| `env`                          | string | Environment (`local`, `dev`, `prod`)          | `"local"`                |
| `storage.driver`               | string | Database driver (`sqlite3`, `postgres`, etc.) | `"sqlite3"`              |
| `storage.dsn`                  | string | Data source name                              | `"storage/scheduler.db"` |
| `storage.limit`                | int    | Pagination limit                              | `10`                     |
| `http_server.address`          | string | Server address                                | `"localhost"`            |
| `http_server.port`             | string | Server port                                   | `"8000"`                 |
| `http_server.timeout`          | string | Read and write timeouts                       | `"4s"`                   |
| `http_server.idle_timeout`     | string | Server idle timeout                           | `"60s"`                  |
| `http_server.file_server_path` | string | Path to static files                          | `"./web"`                |
