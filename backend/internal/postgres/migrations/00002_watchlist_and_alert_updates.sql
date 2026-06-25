-- +goose Up
create table if not exists watchlist_items (
	user_id text not null,
	action_id text not null,
	created_at timestamptz not null,
	primary key (user_id, action_id)
);

insert into watchlist_items (user_id, action_id, created_at)
select user_id, action_id, min(created_at)
from alerts
group by user_id, action_id
on conflict (user_id, action_id) do nothing;

alter table alerts add column if not exists updated_at timestamptz not null default now();

alter table alerts
	add constraint alerts_watchlist_fk
	foreign key (user_id, action_id)
	references watchlist_items (user_id, action_id)
	on delete cascade;

-- +goose Down
alter table alerts drop constraint if exists alerts_watchlist_fk;
alter table alerts drop column if exists updated_at;
drop table if exists watchlist_items;
