CREATE TYPE error_sources_enum AS ENUM ('internal', 'external');
CREATE TABLE pocket_session (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	session_key char(44) NOT NULL UNIQUE,
	session_height integer NOT NULL,
	portal_region_name varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT pk_tbl_0 PRIMARY KEY (id, portal_region_name)
);
CREATE TABLE portal_region (
	portal_region_name varchar NOT NULL,
	CONSTRAINT pk_portal_region PRIMARY KEY (portal_region_name)
);
CREATE TABLE relay (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	pokt_chain_id char(4) NOT NULL,
	endpoint_id varchar NOT NULL,
	session_key char(44) NOT NULL,
	protocol_app_public_key char(64) NOT NULL,
	relay_source_url varchar,
	pokt_node_address char(40),
	pokt_node_domain varchar,
	pokt_node_public_key char(64),
	relay_start_datetime timestamp NOT NULL,
	relay_return_datetime timestamp NOT NULL,
	is_error boolean NOT NULL,
	error_code integer,
	error_name varchar,
	error_message varchar,
	error_source error_sources_enum,
	error_type varchar,
	relay_roundtrip_time float NOT NULL,
	relay_chain_method_ids varchar NOT NULL,
	relay_data_size integer NOT NULL,
	relay_portal_trip_time float NOT NULL,
	relay_node_trip_time float NOT NULL,
	relay_url_is_public_endpoint boolean NOT NULL,
	portal_region_name varchar NOT NULL,
	is_altruist_relay boolean NOT NULL,
	is_user_relay boolean NOT NULL,
	request_id varchar NOT NULL,
	pokt_tx_id varchar,
	gigastake_app_id varchar,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	blocking_plugin varchar,
	CONSTRAINT pk_relay PRIMARY KEY (id, portal_region_name)
);
CREATE TABLE service_record (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	node_public_key char(64) NOT NULL,
	pokt_chain_id char(4) NOT NULL,
	session_key char(44) NOT NULL,
	request_id varchar NOT NULL,
	portal_region_name varchar NOT NULL,
	latency float NOT NULL,
	tickets integer NOT NULL,
	result varchar NOT NULL,
	available boolean NOT NULL,
	successes integer NOT NULL,
	failures integer NOT NULL,
	p90_success_latency float NOT NULL,
	median_success_latency float NOT NULL,
	weighted_success_latency float NOT NULL,
	success_rate float NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	CONSTRAINT pk_service_record PRIMARY KEY (id, portal_region_name)
);
ALTER TABLE relay
ADD CONSTRAINT fk_relay_portal_region FOREIGN KEY (portal_region_name) REFERENCES portal_region(portal_region_name);
ALTER TABLE relay
ADD CONSTRAINT fk_relay_session FOREIGN KEY (session_key) REFERENCES pocket_session(session_key);
ALTER TABLE service_record
ADD CONSTRAINT fk_service_region_portal_region FOREIGN KEY (portal_region_name) REFERENCES portal_region(portal_region_name);
ALTER TABLE service_record
ADD CONSTRAINT fk_service_record_session FOREIGN KEY (session_key) REFERENCES pocket_session(session_key);
ALTER TABLE pocket_session
ADD CONSTRAINT fk_pocket_session_portal_region FOREIGN KEY (portal_region_name) REFERENCES portal_region(portal_region_name);
INSERT INTO portal_region (portal_region_name)
VALUES ('europe-west3'),
	('europe-north1'),
	('europe-west8'),
	('europe-southwest1'),
	('europe-west2'),
	('europe-west9'),
	('us-east4'),
	('us-east5'),
	('us-west2'),
	('us-west1'),
	('northamerica-northeast2'),
	('asia-east2'),
	('asia-northeast1'),
	('asia-northeast3'),
	('asia-south1'),
	('asia-southeast1'),
	('australia-southeast1');
INSERT INTO pocket_session (
		session_key,
		session_height,
		portal_region_name,
		created_at,
		updated_at
	)
VALUES (
		'',
		1,
		'europe-west3',
		TIMESTAMP '1970-01-01 00:00:00',
		TIMESTAMP '1970-01-01 00:00:00'
	);
