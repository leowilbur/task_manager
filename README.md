# Task manager 

## CI status with CircleCI

[![CircleCI](https://circleci.com/gh/leowilbur/task_manager.svg?style=svg)](https://app.circleci.com/pipelines/github/leowilbur/task_manager)
## Project set auto CD with Heroku for testing purpose, just click on icon below: 
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://leowilbur-task-manager.herokuapp.com/swagger.json)

## Development

Use below command to run the service:

```bash
go run main.go
```

You can run the tests:

```bash
go test -v ./...
```
or 
```
make test
```

## Project structure

 - Project using `golang` with [`gin-gonic/gin`](https://github.com/gin-gonic/gin) for rest API,
 - Cron jobs with `github.com/robfig/cron`,
 - `.circleci` contains files related to the CI process,
 - `docs` contains all the documentation packages used for resource embedding,
 - `pkg`:
   - `models` contains all the data model declarations and a basic CRUD layer,
   - `rest` contains an implementation of the REST API,
 - `go.mod` and `go.sum` for package dependencies manage,
 - `task_template.csv` is example `csv` template for import task

## API Document

Postman document with request and example response:

[Click ==> Postman latest API document](https://documenter.getpostman.com/view/8050990/Tz5wVu2j)
[Click ==> Swagger document](https://leowilbur-task-manager.herokuapp.com/swagger.json)

Create Task:
```
curl --location --request POST 'localhost:3000/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Run every minutes",
    "cron": "* * * * *"
}'
```

Delete Task:
```
curl --location --request DELETE 'localhost:3000/tasks/1' 
```

Get Task List
```
curl --location --request GET 'localhost:3000/tasks'
```

Start Tasks:
```
curl --location --request POST 'localhost:3000/task/start' 
```

Stop Tasks:
```
curl --location --request POST 'localhost:3000/task/stop' 
```

Import Tasks:
```
curl --location --request POST 'localhost:3000/task/import' \
--form 'file=@/home/leo/task_template.csv'
```

Export Tasks:
```
curl --location --request POST 'localhost:3000/task/export'
```

## Note
If you have any question or issue, feel free to leave comment.
Thank you