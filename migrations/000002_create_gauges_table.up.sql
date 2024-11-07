CREATE TABLE IF NOT EXISTS gauges (
    id SERIAL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    value DOUBLE PRECISION,
    primary key(name)
);