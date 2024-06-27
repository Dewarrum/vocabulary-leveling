BEGIN;

CREATE TABLE IF NOT EXISTS inits (
    id UUID PRIMARY KEY NOT NULL,
    video_id UUID NOT NULL,
    representation_id TEXT NOT NULL,
    content_location TEXT NOT NULL
);

CREATE UNIQUE INDEX udx_inits_video_id_representation_id ON inits (video_id, representation_id);

COMMIT;