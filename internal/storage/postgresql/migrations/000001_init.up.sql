CREATE TABLE orders
(
    order_uid          VARCHAR PRIMARY KEY,
    track_number       VARCHAR UNIQUE,
    entry              VARCHAR,
    locale             VARCHAR,
    internal_signature VARCHAR,
    customer_id        VARCHAR,
    delivery_service   VARCHAR,
    shardkey           VARCHAR,
    sm_id              INTEGER,
    date_created       TIMESTAMP WITH TIME ZONE,
    oof_shard          VARCHAR
);

CREATE TABLE deliveries
(
    delivery_id SERIAL PRIMARY KEY,
    order_uid   VARCHAR REFERENCES orders (order_uid),
    name        VARCHAR,
    phone       VARCHAR,
    zip         VARCHAR,
    city        VARCHAR,
    address     VARCHAR,
    region      VARCHAR,
    email       VARCHAR
);

CREATE TABLE items
(
    item_id      SERIAL PRIMARY KEY,
    chrt_id      INTEGER,
    track_number VARCHAR REFERENCES orders (track_number),
    price        NUMERIC,
    rid          VARCHAR,
    name         VARCHAR,
    sale         NUMERIC,
    size         VARCHAR,
    total_price  NUMERIC,
    nm_id        INTEGER UNIQUE,
    brand        VARCHAR,
    status       INTEGER
);

CREATE TABLE payments
(
    payment_id    SERIAL PRIMARY KEY,
    order_uid     VARCHAR REFERENCES orders (order_uid),
    request_id    VARCHAR,
    currency      VARCHAR,
    provider      VARCHAR,
    amount        NUMERIC,
    payment_dt    BIGINT,
    bank          VARCHAR,
    delivery_cost NUMERIC,
    goods_total   NUMERIC,
    custom_fee    NUMERIC,
    transaction   VARCHAR UNIQUE
);
