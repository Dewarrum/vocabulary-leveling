BEGIN;

DROP TABLE IF EXISTS chunks;
DROP INDEX IF EXISTS udx_chunks_video_id_sequence;
DROP INDEX IF EXISTS idx_video_id_start_ms_end_ms;

COMMIT;