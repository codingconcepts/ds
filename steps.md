Postgres
``` sql
CREATE TABLE example (
  col_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  col_string VARCHAR(255) NOT NULL,
  col_int INT NOT NULL,
  col_decimal DECIMAL NOT NULL,
  col_date DATE NOT NULL, 
  col_timestamp TIMESTAMPTZ NOT NULL
);

INSERT INTO example (col_string, col_int, col_decimal, col_date, col_timestamp) VALUES
  ('a', 1, 1.0, now()::DATE, now()),
  ('b', 2, 2.0, now()::DATE, now()),
  ('c', 3, 3.0, now()::DATE, now()),
  ('d', 4, 4.0, now()::DATE, now()),
  ('e', 5, 5.0, now()::DATE, now());
```

Cockroach
``` sql
CREATE TABLE example (
  col_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  col_string STRING NOT NULL,
  col_int INT NOT NULL,
  col_decimal DECIMAL NOT NULL,
  col_date DATE NOT NULL, 
  col_timestamp TIMESTAMPTZ NOT NULL
);
```