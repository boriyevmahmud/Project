ALTER TABLE users_test ADD COLUMN IF NOT EXISTS bio text ;
ALTER TABLE users_test ALTER COLUMN email TYPE text ;