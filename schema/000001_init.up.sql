CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(50) NOT NULL,
                       name VARCHAR(50) NOT NULL,
                       username VARCHAR(50) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE typeObject (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL
);

CREATE TABLE objects (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL,
                       server_name VARCHAR(50) NULL,
                       size INT NULL,
                       create_date TIMESTAMP NULL,
                       user_id INT NOT NULL,
                       parent_id INTEGER REFERENCES objects ON DELETE CASCADE,
                       type_object_id INT NOT NULL,
                       FOREIGN KEY (user_id) REFERENCES users(id),
                       FOREIGN KEY (type_object_id) REFERENCES typeObject(id)
);

INSERT INTO typeObject (name) VALUES
                        ('file'),
                        ('directory'),
                        ('main_directory');