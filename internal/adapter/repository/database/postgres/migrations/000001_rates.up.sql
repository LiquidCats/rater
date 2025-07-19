create table if not exists pairs (
    symbol varchar(20) primary key,
    base_asset varchar(5) not null,
    quote_asset varchar(5) not null,
    created_at timestamp not null default current_timestamp
);

insert into pairs (symbol, base_asset, quote_asset) values ('BTC_USD', 'BTC', 'USD');
insert into pairs (symbol, base_asset, quote_asset) values ('BTC_EUR', 'BTC', 'EUR');
insert into pairs (symbol, base_asset, quote_asset) values ('ETH_USD', 'ETH', 'USD');
insert into pairs (symbol, base_asset, quote_asset) values ('ETH_EUR', 'ETH', 'EUR');

create table if not exists providers (
    name varchar(50) primary key,
    created_at timestamp not null default current_timestamp
);

insert into providers (name) values ('cex');
insert into providers (name) values ('coinapi');
insert into providers (name) values ('coingate');
insert into providers (name) values ('coingecko');
insert into providers (name) values ('coinmarketcap');

create table if not exists rates(
    id serial primary key,
    price decimal(10,18) not null,
    pair varchar(20) not null references pairs(symbol),
    provider varchar(50) not null references providers(name),
    ts timestamp not null default current_timestamp,
    created_at timestamp not null default current_timestamp,
    unique(pair, ts)
);

-- create indexes for better query performance
create index if not exists idx_rates_provider on rates(provider);
create index if not exists idx_rates_ts on rates(ts);
