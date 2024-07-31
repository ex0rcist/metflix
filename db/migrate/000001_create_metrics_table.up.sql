CREATE TYPE metricKind AS ENUM ('counter', 'gauge');

CREATE TABLE IF NOT EXISTS metrics(
    id    varchar(255) primary key,
    name  varchar(255) not null,
    kind  metricKind not null,
    value double precision not null
);

CREATE INDEX metrics__idx ON metrics (id);