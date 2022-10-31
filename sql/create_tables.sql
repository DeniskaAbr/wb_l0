create table if not exists orders
(
--     id                 integer generated always as identity
--         constraint orders_id_pkey
--             primary key,
    order_uid          varchar(19) not null
        constraint orders_uid_pkey
            primary key,
    track_number       varchar(14),
    entry              varchar(4),
    locale             varchar(2),
    internal_signature varchar(8),
    customer_id        varchar(4),
    delivery_service   varchar(5),
    shardkey           varchar(2),
    sm_id              integer,
    date_created       varchar(20),
    oof_shard          varchar(2)
);

alter table orders
    owner to market_user;

    create table if not exists items
(
    order_uid     varchar(19) not null
        constraint items_orders_uid_fk
            references orders,
    chrt_id      integer,
    track_number varchar(14),
    price        integer,
    rid          varchar(21),
    name         varchar(255),
    sale         integer,
    size         varchar(3),
    total_price  integer,
    nm_id        integer,
    brand        varchar(255),
    status       integer
);

alter table items
    owner to market_user;

create table if not exists payments
(
    order_uid      varchar(19) not null
        constraint payments_orders_id_fk
            references orders,
    transaction   varchar(20),
    request_id    varchar(20),
    currency      varchar(20),
    provider      varchar(20),
    amount        integer,
    payment_dt    integer,
    bank          varchar(20),
    delivery_cost integer,
    goods_total   integer,
    custom_fee    integer
);

alter table payments
    owner to market_user;

create table if not exists deliveries
(
    order_uid varchar(19) not null
        constraint deliveries_orders_id_fk
            references orders,
    name     varchar(20),
    phone    varchar(20),
    zip      varchar(20),
    city     varchar(20),
    address  varchar(20),
    region   varchar(20),
    email    varchar(30)
);

alter table deliveries
    owner to market_user;


