

/* Create Tables */

CREATE TABLE annotation
(
	id_annotation uuid NOT NULL,
	text varchar(255) NOT NULL,
	id_task uuid NOT NULL,
	id_user uuid NOT NULL,
	PRIMARY KEY (id_annotation)
) WITHOUT OIDS;


CREATE TABLE client_team_xref
(
	id_user uuid NOT NULL,
	id_team uuid NOT NULL
) WITHOUT OIDS;


CREATE TABLE image
(
	id_image uuid NOT NULL,
	width float NOT NULL,
	height float NOT NULL,
	rel_path varchar(255) NOT NULL UNIQUE,
	PRIMARY KEY (id_image)
) WITHOUT OIDS;


CREATE TABLE invitation
(
	id_invitation uuid NOT NULL,
	email varchar(255) NOT NULL,
	expires_at timestamp with time zone NOT NULL,
	is_used boolean NOT NULL,
	id_user uuid NOT NULL,
	PRIMARY KEY (id_invitation)
) WITHOUT OIDS;


CREATE TABLE task
(
	id_task uuid NOT NULL,
	title varchar(255) NOT NULL,
	priority int,
	id_user uuid NOT NULL,
	id_team uuid,
	PRIMARY KEY (id_task)
) WITHOUT OIDS;


CREATE TABLE team
(
	id_team uuid NOT NULL,
	name varchar(255) NOT NULL,
	PRIMARY KEY (id_team)
) WITHOUT OIDS;


CREATE TABLE todo
(
	id_todo uuid NOT NULL,
	text varchar(255) NOT NULL,
	id_complete boolean NOT NULL,
	title varchar(255),
	id_task uuid NOT NULL,
	PRIMARY KEY (id_todo)
) WITHOUT OIDS;


CREATE TABLE users
(
	id_user uuid NOT NULL,
	login varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	password varchar(255) NOT NULL,
	created_at timestamp NOT NULL,
	id_image uuid,
	PRIMARY KEY (id_user)
) WITHOUT OIDS;



/* Create Foreign Keys */

ALTER TABLE users
	ADD FOREIGN KEY (id_image)
	REFERENCES image (id_image)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE annotation
	ADD FOREIGN KEY (id_task)
	REFERENCES task (id_task)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE todo
	ADD FOREIGN KEY (id_task)
	REFERENCES task (id_task)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE client_team_xref
	ADD FOREIGN KEY (id_team)
	REFERENCES team (id_team)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE task
	ADD FOREIGN KEY (id_team)
	REFERENCES team (id_team)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE annotation
	ADD FOREIGN KEY (id_user)
	REFERENCES users (id_user)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE client_team_xref
	ADD FOREIGN KEY (id_user)
	REFERENCES users (id_user)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE invitation
	ADD FOREIGN KEY (id_user)
	REFERENCES users (id_user)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE task
	ADD FOREIGN KEY (id_user)
	REFERENCES users (id_user)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;



/* Comments */

COMMENT ON COLUMN annotation.id_annotation IS 'Суррогатный ключ';
COMMENT ON COLUMN annotation.text IS 'Содержание комментария';
COMMENT ON COLUMN annotation.id_task IS 'Внешний ключ';
COMMENT ON COLUMN annotation.id_user IS 'Внешний ключ';
COMMENT ON COLUMN client_team_xref.id_user IS 'Суррогатный ключ';
COMMENT ON COLUMN client_team_xref.id_team IS 'Суррогатный ключ';
COMMENT ON COLUMN image.id_image IS 'Суррогатный ключ';
COMMENT ON COLUMN invitation.id_invitation IS 'Суррогатный ключ';
COMMENT ON COLUMN invitation.email IS 'Электронаая почта приглашённого';
COMMENT ON COLUMN invitation.expires_at IS 'Срок истечения приглашения';
COMMENT ON COLUMN invitation.is_used IS 'Использовано ли приглашение';
COMMENT ON COLUMN invitation.id_user IS 'Суррогатный ключ';
COMMENT ON COLUMN task.id_task IS 'Суррогатный ключ';
COMMENT ON COLUMN task.title IS 'Название задачи';
COMMENT ON COLUMN task.priority IS 'Приоритет расстановки задач в интерфейсе';
COMMENT ON COLUMN task.id_user IS 'Суррогатный ключ';
COMMENT ON COLUMN task.id_team IS 'Суррогатный ключ';
COMMENT ON COLUMN team.id_team IS 'Суррогатный ключ';
COMMENT ON COLUMN team.name IS 'Название группы';
COMMENT ON COLUMN todo.id_todo IS 'Суррогатный ключ';
COMMENT ON COLUMN todo.text IS 'Текст подзадачи';
COMMENT ON COLUMN todo.id_complete IS 'Выполнена ли подзадача';
COMMENT ON COLUMN todo.title IS 'Заголовок подзадачи';
COMMENT ON COLUMN todo.id_task IS 'Внешний ключ';
COMMENT ON COLUMN users.id_user IS 'Суррогатный ключ';
COMMENT ON COLUMN users.login IS 'Логин пользователя';
COMMENT ON COLUMN users.email IS 'Электронная почта пользователя';
COMMENT ON COLUMN users.password IS 'Пароль пользователя';
COMMENT ON COLUMN users.id_image IS 'Внешний ключ';



