create type gender AS ENUM('male', 'female');

create table users(
    id serial primary key,
    username varchar(20),
    password varchar(20),
    age int,
    gender gender
);

insert into users (username, password, age, gender) VALUES
    ('lixia', '1123', 20, 'female');

-- 查询
SELECT *
FROM "users"
WHERE (username = 'lixia');
