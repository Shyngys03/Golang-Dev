ALTER TABLE drills ADD CONSTRAINT weight_check CHECK (weight > 0);
ALTER TABLE drills ADD CONSTRAINT drills_length_check CHECK (cable_length > 0);
ALTER TABLE drills ADD CONSTRAINT chuck_diameter_check CHECK (chuck_diameter > 0 and chuck_diameter < 30);