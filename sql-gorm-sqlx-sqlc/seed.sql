INSERT INTO users (id, name, email) VALUES (1000, "Alice", "alice@domain.com");
INSERT INTO users (id, name, email) VALUES (1001, "Bob", "bob@domain.com");
INSERT INTO users (id, name, email) VALUES (1002, "Charlie", "chalie@doman.com");

INSERT INTO blogs (id, name, url) VALUES (2000, "devsecops", "http://devsecops.com");
INSERT INTO blogs (id, name, url) VALUES (2001, "devops", "http://devops.com");
INSERT INTO blogs (id, name, url) VALUES (2002, "oop", "http://oop.com");

INSERT INTO posts (title, content, user_id, blog_id) VALUES ("", "", 1000, 2000);
INSERT INTO posts (title, content, user_id, blog_id) VALUES ("", "", 1001, 2001);
INSERT INTO posts (title, content, user_id, blog_id) VALUES ("", "", 1002, 2002);
INSERT INTO posts (title, content, user_id, blog_id) VALUES ("", "", 1000, 2001);
INSERT INTO posts (title, content, user_id, blog_id) VALUES ("", "", 1001, 2002);
