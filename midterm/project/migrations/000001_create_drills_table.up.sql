CREATE TABLE IF NOT EXISTS "drills" (
    "id" bigserial PRIMARY KEY,
    "weight" double precision NOT NULL,
    "name" varchar NOT NULL,
    "cable_length" double precision NOT NULL,
    "work_time" integer NOT NULL,
    "chuck_diameter" integer not NULL
);
