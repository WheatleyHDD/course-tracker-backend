CREATE VIEW user_info AS
SELECT u.email, u.first_name, u.second_name, u.middle_name, u.perms
FROM users u