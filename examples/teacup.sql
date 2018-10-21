-- Automated with API

-- Table creation.
CREATE TABLE IF NOT EXISTS posts (
  uid              serial PRIMARY KEY,
  optional_summary varchar(192),
  title            varchar(128) UNIQUE       NOT NULL,
  body             text                      NOT NULL,
  post_date        date DEFAULT CURRENT_DATE NOT NULL
) WITH (OIDS = FALSE
);

-- Summary function for posts.
CREATE OR REPLACE FUNCTION summary(rec posts)
  RETURNS varchar(192)
LANGUAGE SQL
AS
$$
SELECT CASE
  WHEN $1.optional_summary IS NULL
    THEN CASE
      WHEN length($1.body) > 192
        THEN $1.body :: varchar(191) || '…'
      ELSE $1.body :: varchar(192)
    END
  ELSE $1.optional_summary
END
$$;

------------------------------------------------------------------------------------------------------------------------
-- Example manual insertions.

-- Quick and dirty table dropping. `cascade` is required due to the summary function.
drop table posts cascade;

-- Insertions into `posts`
INSERT INTO posts (title, body)
VALUES ('Example 1', 'Here''s some content posted today.');
-- with explicit date
INSERT INTO posts (title, body, post_date)
VALUES ('Example 1', 'Here''s some content posted today.', '2018-10-25');