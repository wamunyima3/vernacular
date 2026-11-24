-- materialise_release copies a built release into its own schema, so pinned
-- reads run the live queries against frozen rows through search_path (ADR-0012).
-- It runs in the release build's transaction, after release_items is filled.
-- The caller names the schema (releases.Schema), so Go holds the one naming rule.

-- +goose Up
-- trust.release looks up the releases an item is in.
CREATE INDEX release_items_item_id ON release_items (item_id);

-- +goose StatementBegin
CREATE FUNCTION materialise_release(release_id uuid, s text) RETURNS void
LANGUAGE plpgsql AS $$
DECLARE
    -- Each copied table, with the column that ties its rows to a released item.
    copies constant text[][] := ARRAY[
        ['entries', 'id'], ['senses', 'id'], ['examples', 'id'], ['phrases', 'id'],
        ['variants', 'entry_id'], ['item_sources', 'item_id']];
    tables text[] := ARRAY['items'];
    c      text[];
    t      text;
    ix     record;
BEGIN
    EXECUTE format('CREATE SCHEMA %I', s);
    EXECUTE format('CREATE TABLE %I.items AS SELECT i.* FROM public.items i
                    JOIN public.release_items r ON r.item_id = i.id WHERE r.release_id = %L', s, release_id);
    FOREACH c SLICE 1 IN ARRAY copies LOOP
        EXECUTE format('CREATE TABLE %I.%I AS SELECT x.* FROM public.%I x
                        WHERE x.%I IN (SELECT id FROM %I.items)', s, c[1], c[1], c[2], s);
        tables := tables || c[1];
    END LOOP;
    -- The live indexes, so pinned search and lookup stay fast, and the primary
    -- keys, which queries grouping by id rely on.
    FOR ix IN
        SELECT x.tablename, x.indexname, x.indexdef, i.indisprimary
        FROM pg_indexes x JOIN pg_index i ON i.indexrelid = format('public.%I', x.indexname)::regclass
        WHERE x.schemaname = 'public' AND x.tablename = ANY(tables)
    LOOP
        EXECUTE replace(ix.indexdef, ' ON public.' || ix.tablename || ' USING',
                        format(' ON %I.%I USING', s, ix.tablename));
        IF ix.indisprimary THEN
            EXECUTE format('ALTER TABLE %I.%I ADD PRIMARY KEY USING INDEX %I', s, ix.tablename, ix.indexname);
        END IF;
    END LOOP;
    -- Planner statistics now, rather than when autovacuum gets to the tables.
    FOREACH t IN ARRAY tables LOOP
        EXECUTE format('ANALYZE %I.%I', s, t);
    END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION materialise_release(uuid, text);
DROP INDEX release_items_item_id;
