create table if not exists drills(
    id bigserial primary key, 
    weigth float,
    name varchar(20),
    cable_length float,
    worktime int,
    diameter int
);