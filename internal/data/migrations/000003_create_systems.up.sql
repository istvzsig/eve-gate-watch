CREATE TABLE systems (
	id BIGINT PRIMARY KEY,
	name TEXT NOT NULL,
	constellation_id BIGINT NOT NULL,
	region_id BIGINT NOT NULL,
	security_status DOUBLE PRECISION NOT NULL,
	wormhole_class_id INTEGER,
	CONSTRAINT fk_systems_constellation FOREIGN KEY (constellation_id) REFERENCES constellations (id) ON DELETE RESTRICT,
	CONSTRAINT fk_systems_region FOREIGN KEY (region_id) REFERENCES regions (id) ON DELETE RESTRICT
);

CREATE UNIQUE INDEX idx_systems_name ON systems (name);

CREATE INDEX idx_systems_constellation_id ON systems (constellation_id);

CREATE INDEX idx_systems_region_id ON systems (region_id);

CREATE INDEX idx_systems_wormhole_class_id ON systems (wormhole_class_id);