CREATE DATABASE management;
USE management;

CREATE TABLE managers(
	manager_id INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(15),
    surname VARCHAR(15),
    email VARCHAR(25) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
	PRIMARY KEY(manager_id)
);

CREATE TABLE employees(
	employee_id INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(15),
    surname VARCHAR(15),
    email VARCHAR (25) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    manager_id INT NOT NULL,
    salary DOUBLE,
    bonuses DOUBLE DEFAULT 0,
    PRIMARY KEY(employee_id),
    FOREIGN KEY(manager_id) REFERENCES managers(manager_id)
);

CREATE TABLE requests(
	request_id INT NOT NULL AUTO_INCREMENT,
    request_type VARCHAR(15),
    request_body VARCHAR(100),
    request_date DATE,
    approval VARCHAR(3) DEFAULT 'No',
    emp_id INT NOT NULL,
    manager_id INT NOT NULL,
    PRIMARY KEY(request_id),
    FOREIGN KEY(emp_id) REFERENCES employees(employee_id),
    FOREIGN KEY(manager_id) REFERENCES managers(manager_id)
);

CREATE TABLE tasks(
	task_id INT NOT NULL AUTO_INCREMENT,
    taskt_title VARCHAR(20),
    task_body VARCHAR(100),
    date_created DATE,
    due_date DATE,
    bonus DOUBLE DEFAULT 0,
    emp_id INT NOT NULL,
    manager_id INT NOT NULL,
    PRIMARY KEY(task_id),
    FOREIGN KEY(emp_id) REFERENCES employees(employee_id),
    FOREIGN KEY(manager_id) REFERENCES managers(manager_id)
);


SHOW CREATE TABLE tasks;
-- drop not null ma heret
ALTER TABLE tasks
ALTER COLUMN emp_id DROP DEFAULT;


alter table tasks
modify column emp_id INT Default 0;

ALTER TABLE tasks MODIFY COLUMN emp_id int;

ALTER TABLE tasks MODIFY COLUMN manager_id int NOT NULL;
-- mos e prekt tasks mo

ALTER TABLE requests MODIFY manager_id INT;
ALTER TABLE requests
ALTER COLUMN manager_id DROP DEFAULT;

ALTER TABLE requests DROP FOREIGN KEY requests_ibfk_2;
ALTER TABLE requests DROP COLUMN manager_id;

ALTER TABLE requests MODIFY COLUMN approval bool DEFAULT 0;
SHOW CREATE table requests;

CREATE TABLE app_requests(
	request_id INT NOT NULL AUTO_INCREMENT,
    request_type VARCHAR(15),
    request_body VARCHAR(100),
    request_date DATE,
    emp_id INT NOT NULL,
    PRIMARY KEY(request_id),
    FOREIGN KEY(emp_id) REFERENCES employees(employee_id)
);

DROP TABLE app_requests;

ALTER table requests modify column emp_id INT;

SHOW CREATE TABLE requests;

ALTER TABLE `requests`
	DROP FOREIGN KEY `requests_ibfk_2`;

alter table requests
  ADD constraint requests_ibfk_1
  foreign key (emp_id)
  references employees (employee_id) on delete cascade;