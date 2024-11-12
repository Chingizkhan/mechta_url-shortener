begin;

create table if not exists "url"(
    link text primary key,
    source_url text not null,
    visits integer default 0,
    created_at timestamp(6) default now(),
    updated_at timestamp(6) default now(),
    expire_at timestamp(6) default now()+ interval '30 days'
);

end;