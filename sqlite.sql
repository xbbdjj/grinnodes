CREATE TABLE IF NOT EXISTS "peer" (
	 "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "ip" text NOT NULL DEFAULT '',
	 "port" integer NOT NULL DEFAULT 0,
	 "node_user_agent" text NOT NULL DEFAULT '',
	 "node_protocol_version" integer NOT NULL DEFAULT 0,
	 "node_capabilities" integer NOT NULL DEFAULT 0,
	 "node_height" integer NOT NULL DEFAULT 0,
	 "node_total_difficulty" integer NOT NULL DEFAULT 0,
	 "genesis" text NOT NULL DEFAULT '',
	 "p2p_first_seen" integer NOT NULL DEFAULT 0,
	 "p2p_last_seen" integer NOT NULL DEFAULT 0,
	 "p2p_first_connected" integer NOT NULL DEFAULT 0,
	 "p2p_last_connected" integer NOT NULL DEFAULT 0,
	 "p2p_failed" integer NOT NULL DEFAULT 0,
	 "api_first_seen" integer NOT NULL DEFAULT 0,
	 "api_last_seen" integer NOT NULL DEFAULT 0,
	 "ip_state" integer NOT NULL DEFAULT 0,
	 "ip_continent_code" text NOT NULL DEFAULT '',
	 "ip_country_code" text NOT NULL DEFAULT '',
	 "ip_country_name" text NOT NULL DEFAULT '',
	 "ip_city" text NOT NULL DEFAULT '',
	 "ip_latitude" text NOT NULL DEFAULT '',
	 "ip_longitude" text NOT NULL DEFAULT '',
	 "ip_rdns" text NOT NULL DEFAULT '',
	 "ip_asn" text NOT NULL DEFAULT '',
	 "ip_org" text NOT NULL DEFAULT '',
	CONSTRAINT "ip" UNIQUE ("ip" ASC) ON CONFLICT ABORT
);


CREATE TABLE IF NOT EXISTS "chart" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "time" integer NOT NULL DEFAULT 0,
  "peer_total" integer NOT NULL DEFAULT 0,
  "peer_public_total" integer NOT NULL DEFAULT 0
);