CREATE INDEX IF NOT EXISTS drill_weight_idx ON drills USING GIN (to_tsvector('simple', weight));
CREATE INDEX IF NOT EXISTS drill_name_idx ON drills USING GIN (to_tsvector('simple', name));