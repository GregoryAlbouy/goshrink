package database

const userSchema = `
CREATE TABLE IF NOT EXISTS user (
	id INT UNSIGNED AUTO_INCREMENT,
	username VARCHAR(20) UNIQUE,
	password VARCHAR(255),
	email VARCHAR(255) UNIQUE,
	PRIMARY KEY (id)
)
`

const avatarSchema = `
CREATE TABLE IF NOT EXISTS avatar (
	id INT UNSIGNED AUTO_INCREMENT,
	user_id INT UNSIGNED UNIQUE,
	avatar_url VARCHAR(255),
	PRIMARY KEY (id),
	CONSTRAINT fk_user_id
		FOREIGN KEY (user_id)
		REFERENCES user(id)
		ON DELETE CASCADE
)`

const userAvatarSchema = `
CREATE OR REPLACE VIEW V_user_avatar
AS SELECT u.id, u.username, u.email, COALESCE(a.avatar_url, '') AS avatar_url
FROM user u
LEFT JOIN avatar a
	ON u.id = a.user_id
`
