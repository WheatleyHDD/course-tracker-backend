CREATE TABLE course_applications (
  id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  course_name varchar(255) NOT NULL,
  student varchar(255) NOT NULL,
  cost integer NOT NULL,
  start_date date NOT NULL,
  end_date date NOT NULL,
  point varchar(255) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE users (
  id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  first_name varchar(255) NOT NULL,
  second_name varchar(255) NOT NULL,
  middle_name varchar(255) NOT NULL,
  perms integer NOT NULL,
  PRIMARY KEY (email)
);

CREATE TABLE comms (
  id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  application_id bigint NOT NULL,
  sender varchar(255) NOT NULL,
  comm_timestamp timestamp NOT NULL,
  message_text varchar(255) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT comments_application_id_application_id_foreign FOREIGN KEY (application_id) REFERENCES course_applications (id),
  CONSTRAINT comments_sender_users_id_foreign FOREIGN KEY (sender) REFERENCES users (email)
);

CREATE TABLE course (
  id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  application_id bigint NOT NULL,
  tutor_id varchar(255) NOT NULL,
  department varchar(255) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT course_tutor_id_users_id_foreign FOREIGN KEY (tutor_id) REFERENCES users (email),
  CONSTRAINT course_application_id_course_applications_id_foreign FOREIGN KEY (application_id) REFERENCES course_applications (id)
);

CREATE TABLE statuses (
  id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  application_id bigint NOT NULL,
  changer varchar(255) NOT NULL,
  change_date timestamp NOT NULL,
  status integer NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT statuses_application_id_application_id_foreign FOREIGN KEY (application_id) REFERENCES course_applications (id),
  CONSTRAINT statuses_changer_users_id_foreign FOREIGN KEY (changer) REFERENCES users (email)
);
COMMENT ON COLUMN statuses.status IS 'status:
0 - На согласовании
1 - Отклонено
2 - Новая (На рассмотрении админом)
3 - В работе
4 - Согласование документов (Внутренний)
5 - Оплата (Внутренний)
6 - Ожидает обучения
7 - Пройдено';

CREATE TABLE tokens (
  users_id varchar(255) NOT NULL,
  access_token varchar(255) NOT NULL,
  CONSTRAINT tokens_users_id_users_id_foreign FOREIGN KEY (users_id) REFERENCES users (email)
);

CREATE VIEW cources_and_statuses AS
SELECT apps.id, apps.course_name, apps.student, apps.cost, apps.start_date, apps.end_date, apps.point, s.status, s.changer, s.change_date
FROM course_applications apps
LEFT JOIN (
	SELECT sd.application_id, max(sd.id) id
	FROM statuses sd
	GROUP BY sd.application_id
) sids ON sids.application_id = apps.id
LEFT JOIN statuses s ON s.id = sids.id;

CREATE VIEW user_info AS
SELECT u.email, u.first_name, u.second_name, u.middle_name, u.perms
FROM users u