CREATE TABLE events (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
    "name" TEXT NOT NULL,
    "date" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "description" TEXT NOT NULL,
    "user_id" integer NOT NULL
);


CREATE TRIGGER update_event_updated_at_timestamp
AFTER UPDATE ON events
FOR EACH ROW
BEGIN
    UPDATE events SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
