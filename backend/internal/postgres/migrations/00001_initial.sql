-- +goose Up
create table if not exists actions (
	id text primary key,
	symbol text not null unique,
	name text not null,
	exchange text not null default '',
	active boolean not null default true,
	created_at timestamptz not null,
	updated_at timestamptz not null
);

insert into actions (id, symbol, name, exchange, active, created_at, updated_at)
values
	('action-petr4', 'PETR4', 'Petrobras PN', 'B3', true, now(), now()),
	('action-vale3', 'VALE3', 'Vale ON', 'B3', true, now(), now()),
	('action-bbsa3', 'BBAS3', 'Banco do Brasil ON', 'B3', true, now(), now())
on conflict (symbol) do nothing;

create table if not exists alerts (
	id text primary key,
	user_id text not null,
	action_id text not null,
	symbol text not null,
	action_name text not null default '',
	target_price double precision not null,
	direction text not null,
	device_token text not null default '',
	status text not null,
	created_at timestamptz not null,
	triggered_at timestamptz null
);

alter table alerts add column if not exists action_id text not null default '';
alter table alerts add column if not exists action_name text not null default '';

update alerts
set action_id = actions.id,
	action_name = actions.name
from actions
where alerts.action_id = ''
	and upper(alerts.symbol) = upper(actions.symbol);

create index if not exists alerts_user_id_idx on alerts (user_id);
create index if not exists alerts_symbol_status_idx on alerts (symbol, status);
create index if not exists alerts_action_id_idx on alerts (action_id);

create table if not exists device_registrations (
	user_id text primary key,
	device_token text not null,
	platform text not null default '',
	created_at timestamptz not null
);

-- +goose Down
drop table if exists device_registrations;
drop index if exists alerts_action_id_idx;
drop index if exists alerts_symbol_status_idx;
drop index if exists alerts_user_id_idx;
drop table if exists alerts;
drop table if exists actions;
