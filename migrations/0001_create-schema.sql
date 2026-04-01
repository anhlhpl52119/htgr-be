-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS departments (
  id INT PRIMARY KEY,
  code VARCHAR(30) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
);


INSERT INTO
  departments (id, code, name)
VALUES
  (100, 'DEV_FE', 'Front end Dev Team'),
  (200, 'DEV_BE', 'Back end Dev Team'),
  (300, 'BO_HR', 'Back Office HR'),
  (300, 'BO_FI', 'Back Office Finance');


CREATE TABLE IF NOT EXISTS roles (
  id INT PRIMARY KEY,
  code VARCHAR(30) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
  description TEXT
);


INSERT INTO
  roles (id, code, name)
VALUES
  (100, 'EMPLOYEE', 'Employee'),
  (200, 'HR', 'HR'),
  (300, 'MANAGER', 'Manager'),
  (999, 'ADMIN', 'Admin');


CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
);


CREATE TABLE user_roles (
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  role_id INT REFERENCES roles (id) DEFAULT 100,
  PRIMARY KEY (user_id, role_id)
);


CREATE TABLE IF NOT EXISTS employees (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users (id) ON DELETE SET NULL,
  employee_code VARCHAR(50) NOT NULL UNIQUE,
  work_email VARCHAR(100) NOT NULL UNIQUE,
  full_name VARCHAR(100) NOT NULL,
  alias_name VARCHAR(50),
  avatar_img_url VARCHAR(255),
  full_name_normalized VARCHAR(100) NOT NULL,
  phone_number VARCHAR(20),
  manager_id UUID REFERENCES employees (id) ON DELETE SET NULL,
  department_id INT REFERENCES departments (id),
  join_date TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;


DROP TABLE IF EXISTS user_roles;


DROP TABLE IF EXISTS roles;


DROP TABLE IF EXISTS departments;


DROP TABLE IF EXISTS employees;


-- +goose StatementEnd
