CREATE TABLE media (
    id  TEXT PRIMARY KEY,
    sha256 TEXT NOT NULL UNIQUE,
    kind TEXT check(kind = 'photo' or kind = 'video') NOT NULL,
    original_name TEXT NOT NULL,
    original_path TEXT NOT NULL UNIQUE,
    orientation INTEGER NOT NULL DEFAULT 0,
    size_bytes INTEGER NOT NULL,
    mime_type TEXT NOT NULL,
    taken_at TEXT,
    width INTEGER,
    height INTEGER,
    thumb_state TEXT NOT NULL DEFAULT 'pending',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
CREATE INDEX idx_photos_taken_at ON media(taken_at DESC, id DESC);
CREATE INDEX idx_thumbs ON media(thumb_state) WHERE thumb_state != 'done';
