CREATE TABLE IF NOT EXISTS counters (
    id SERIAL,
    name TEXT NOT NULL,
	type TEXT NOT NULL,
    value bigint,
    primary key(name)
);