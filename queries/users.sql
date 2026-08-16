<-- GetById one
select * from users where id = $id;

<-- GetList many
select * from users where (name ilike $data.name or surname ilike $data.surname) order by id desc limit $limit;

<-- Count count
select count(*) from users where (name ilike $data.name or surname ilike $data.surname);

<-- Create exec
insert into users ...;