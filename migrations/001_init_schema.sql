-- Add uuid support for the reports table
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

create table citizen_reports (
    id uuid primary key default gen_random_uuid(),
    region_id bigint not null,
    district_id bigint not null,
    infrastructure_name varchar(255) not null,
    sector_id bigint not null,
    description text,
    photo_path varchar(500),
    created_at timestamp without time zone default now(),
    updated_at timestamp without time zone
);


create or replace function set_updated_at()
returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;

create trigger trg_set_updated_at
before insert or update on citizen_reports
for each row
execute function set_updated_at();