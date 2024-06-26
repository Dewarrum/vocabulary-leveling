BEGIN;

CREATE TABLE IF NOT EXISTS manifests (
    id UUID PRIMARY KEY NOT NULL,
    video_id UUID NOT NULL,
    meta JSONB NOT NULL
);

CREATE UNIQUE INDEX udx_manifests_video_id ON manifests (video_id);

COMMIT;