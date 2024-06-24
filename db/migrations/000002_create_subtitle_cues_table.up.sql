BEGIN;

CREATE TABLE IF NOT EXISTS subtitle_cues (
    id UUID PRIMARY KEY NOT NULL,
    video_id UUID NOT NULL,
    sequence INTEGER NOT NULL,
    start_ms BIGINT NOT NULL,
    end_ms BIGINT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX udx_subtitle_cues_video_id_sequence ON subtitle_cues (video_id, sequence);

COMMIT;