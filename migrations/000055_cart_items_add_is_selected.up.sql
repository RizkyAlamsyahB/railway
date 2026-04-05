ALTER TABLE cart_items
ADD COLUMN is_selected BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE cart_items
SET is_selected = TRUE
WHERE is_selected IS DISTINCT FROM TRUE;
