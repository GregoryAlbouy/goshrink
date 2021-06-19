# Database

The package `database` implements all services related to data storing.

## Table user

| Name     | Type         | Settings    | References |
| -------- | ------------ | ----------- | ---------- |
| id       | INT          | PRIMARY KEY |            |
| username | VARCHAR(20)  |             |            |
| password | VARCHAR(255) |             |            |
| email    | VARCHAR(255) |             |            |

## Table avatar

| Name       | Type         | Settings           | References |
| ---------- | ------------ | ------------------ | ---------- |
| id         | INT          | PRIMARY KEY        |            |
| user_id    | INT          | UNIQUE FOREIGN KEY | user.id    |
| avatar_url | VARCHAR(255) |                    |            |
