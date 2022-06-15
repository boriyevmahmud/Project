ALTER TABLE users_test DROP COLUMN IF EXISTS bio text NOT null;
ALTER TABLE users_test ALTER COLUMN email TYPE varchar(255) NOT null;