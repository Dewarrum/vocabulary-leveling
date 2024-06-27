BEGIN;

CREATE TABLE IF NOT EXISTS chunks (
    id UUID PRIMARY KEY NOT NULL,
    video_id UUID NOT NULL,
    representation_id TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    content_location TEXT NOT NULL,
    start_ms BIGINT NOT NULL,
    end_ms BIGINT NOT NULL
);

CREATE UNIQUE INDEX udx_chunks_video_id_representation_id_sequence ON chunks (video_id, representation_id, sequence);
CREATE INDEX idx_video_id_start_ms_end_ms ON chunks (video_id, start_ms, end_ms);

COMMIT;