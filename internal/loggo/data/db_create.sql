PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS entry
(
    entry_id  INTEGER   NOT NULL PRIMARY KEY,
    stored_at TIMESTAMP NOT NULL DEFAULT (DATETIME('NOW')),
    msg       TEXT      NOT NULL
);
CREATE INDEX IF NOT EXISTS index_entry_stored_at ON entry (stored_at);

CREATE TABLE IF NOT EXISTS meta
(
    entry_id INTEGER NOT NULL,
    tag      TEXT    NOT NULL,
    value            NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (entry_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS index_meta_tag ON meta (tag);

DROP VIEW IF EXISTS entry2;
CREATE VIEW IF NOT EXISTS entry2 AS
    WITH q AS (
        SELECT entry.entry_id, stored_at, msg, tag, value
        FROM entry
                 LEFT OUTER JOIN meta USING (entry_id)
        UNION
        SELECT entry.entry_id, stored_at, msg, tag, value
        FROM meta
                 LEFT OUTER JOIN entry USING (entry_id)
    )
    SELECT entry_id, stored_at,
           IFNULL(msg || ' ' || GROUP_CONCAT(tag || '=' || value, ' '), msg) AS msg
    FROM q
    GROUP BY entry_id;;

-- SELECT entry_id, stored_at, msg
-- FROM entry2
--          INNER JOIN meta USING(entry_id)
-- WHERE tag = 'addr';

