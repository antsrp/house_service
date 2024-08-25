CREATE OR REPLACE FUNCTION update_house()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    UPDATE houses
    SET updated_at = now()
    WHERE id = new.house_id;
    return new;
END $$;

CREATE OR REPLACE TRIGGER house_upd_on_flat_upd
AFTER UPDATE ON flats
FOR EACH ROW
EXECUTE FUNCTION update_house();

CREATE OR REPLACE TRIGGER house_upd_on_flat_create
AFTER INSERT ON flats
FOR EACH ROW
EXECUTE FUNCTION update_house();