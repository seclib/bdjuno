CREATE TABLE btc_header_info
(
    height BIGINT UNIQUE PRIMARY KEY,
    hash   TEXT NOT NULL UNIQUE,
    header TEXT NOT NULL UNIQUE,
    work   TEXT NOT NULL
);
CREATE INDEX btc_header_info_height_index ON btc_header_info (height);
CREATE INDEX btc_header_info_hash_index ON btc_header_info (hash);
CREATE INDEX btc_header_info_header_index ON btc_header_info (header);
ALTER TABLE btc_header_info
    SET (
    autovacuum_vacuum_scale_factor = 0,
    autovacuum_analyze_scale_factor = 0,
    autovacuum_vacuum_threshold = 10000,
    autovacuum_analyze_threshold = 10000
    );
