CREATE TABLE constellations (
	id BIGINT PRIMARY KEY,
	name TEXT NOT NULL,
	region_id BIGINT NOT NULL,
	CONSTRAINT fk_constellations_region FOREIGN KEY (region_id) REFERENCES regions (id) ON DELETE RESTRICT
);

CREATE INDEX idx_constellations_region_id ON constellations (region_id);

CREATE UNIQUE INDEX idx_constellations_region_name ON constellations (region_id, name);