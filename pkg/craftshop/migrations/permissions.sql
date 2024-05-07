-- Set the activated field for alice@example.com to true.
UPDATE users SET activated = true WHERE email = 'xahn2nkw@example.com';


-- Give all members the 'movies:read' permission
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'products:read') FROM users;
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'sellers:read') FROM users;
select * from users_permissions;
TRUNCATE TABLE users_permissions;

-- Give faith@example.com the 'movies:write' permission
INSERT INTO users_permissions
    VALUES (
    (SELECT id FROM users WHERE email = 'uz9xtzlu@example.com'),
    (SELECT id FROM permissions WHERE code = 'products:write')
);

INSERT INTO users_permissions
    VALUES (
    (SELECT id FROM users WHERE email = 'nastya@example.com'),
    (SELECT id FROM permissions WHERE code = 'sellers:write')
);


-- List all activated members and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;

