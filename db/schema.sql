DROP TABLE IF EXISTS todo;
CREATE TABLE IF NOT EXISTS todo (
    id INT AUTO_INCREMENT PRIMARY KEY,
    channel_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_at TIMESTAMP NOT NULL,
    owner_id INT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES user(id)
);

CREATE TABLE IF NOT EXISTS user (
    id INT PRIMARY KEY,
    traq_id VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_todo_relation (
    user_id INT NOT NULL,
    todo_id INT NOT NULL,
    PRIMARY KEY (user_id, todo_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (todo_id) REFERENCES todo(id)
);
