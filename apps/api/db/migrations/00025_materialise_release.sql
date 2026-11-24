-- materialise_release copies a built release into its own schema, so pinned
-- reads run the live queries against frozen rows through search_path (ADR-0012).
-- It runs in the release build's transaction, after release_items is filled.
-- The caller names the schema (releases.Schema), so Go holds the one naming rule.

-- +goose Up
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
    -- The live indexes, so pinned search and lookup stay fast.
    FOR ix IN
        SELECT tablename, indexdef FROM pg_indexes
        WHERE schemaname = 'public' AND tablename = ANY(tables)
    LOOP
        EXECUTE replace(ix.indexdef, ' ON public.' || ix.tablename || ' USING',
                        format(' ON %I.%I USING', s, ix.tablename));
    END LOOP;
    -- Planner statistics now, rather than when autovacuum gets to the tables.
    FOREACH t IN ARRAY tables LOOP
        EXECUTE format('ANALYZE %I.%I', s, t);
    END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION materialise_release(uuid, text);
