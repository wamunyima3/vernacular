-- Supported and planned languages, in rollout order (requirements §5).

-- +goose Up
CREATE TABLE languages (
    code       text        PRIMARY KEY CHECK (code ~ '^[a-z]{3}$'),
    name       text        NOT NULL,
    own_name   text        NOT NULL,
    wave       smallint    NOT NULL CHECK (wave >= 1),
    status     text        NOT NULL DEFAULT 'planned' CHECK (status IN ('active', 'planned')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO languages (code, name, own_name, wave, status) VALUES
    ('bem', 'Bemba',  'Ichibemba',  1, 'active'),
    ('nya', 'Nyanja', 'Chinyanja',  1, 'active'),
    ('toi', 'Tonga',  'Chitonga',   2, 'planned'),
    ('loz', 'Lozi',   'Silozi',     2, 'planned'),
    ('kqn', 'Kaonde', 'Kikaonde',   3, 'planned'),
    ('lun', 'Lunda',  'Chilunda',   3, 'planned'),
    ('lue', 'Luvale', 'Luvale',     3, 'planned');

-- +goose Down
DROP TABLE languages;
