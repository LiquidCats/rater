create table if not exists rates(
    id serial primary key,
    price decimal(10,18) not null,
    pair varchar(10) not null,
    provider varchar(20) not null,
    ts timestamp not null default current_timestamp,
    created_at timestamp not null default current_timestamp,
    UNIQUE(provider, pair, created_at)
);

-- create indexes for better query performance
create index if not exists idx_bitcoin_prices_provider_pair on rates(provider, pair);
create index if not exists idx_bitcoin_prices_created_at on rates(created_at);