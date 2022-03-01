CREATE TABLE IF NOT EXISTS servicetable (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    description   TEXT,
    versions      JSON
);

CREATE EXTENSION pg_trgm;
CREATE INDEX trgm_idx on servicetable using gin ((name || ' ' || description) gin_trgm_ops);

INSERT INTO servicetable (id, name, description, versions)
VALUES (0, 'Notifications', 'Customizable notifications', '[{ "semver": "0.0.0" }, { "semver": "0.1.0" }, { "semver": "2.1.0" }]');

INSERT INTO servicetable (id, name, description, versions)
VALUES (1, 'Notifications', 'Cloud storage monitoring', '[{ "semver": "0.0.0" }, { "semver": "0.1.0" }, { "semver": "2.1.0" }]');

INSERT INTO servicetable (id, name, description, versions)
VALUES (2, 'Reporting', 'Customizable periodic reports on service and account performace', '[{ "semver": "0.0.0" }, { "semver": "0.1.0" }, { "semver": "2.1.0" }]');

INSERT INTO servicetable (id, name, description)
VALUES (3, 'Security', 'Keep it safe, keep it secret');

INSERT INTO servicetable (id, name, description, versions)
VALUES (4, 'Contact Us', 'Multi-channel communication tools', '[{ "semver": "0.0.0" }]');

INSERT INTO servicetable (id, name)
VALUES (5, 'Contact Us');

INSERT INTO servicetable (id, name, description)
VALUES (6, 'Contact Us', 'Suite of notification tools');

INSERT INTO servicetable (id, name, description)
VALUES (7, 'Locate Us', 'Geolocation publishing service');

INSERT INTO servicetable (id, name)
VALUES (8, 'Locate Us');

INSERT INTO servicetable (id, name, description, versions)
VALUES (9, 'Cloud Functions', 'Execute lightweight, scalable program snippets', '[{ "semver": "0.1.0" }, { "semver": "2.2.0" }]');

INSERT INTO servicetable (id, name)
VALUES (10, 'Oauth');

INSERT INTO servicetable (id, name, description)
VALUES (11, 'Oauth', 'Securely integrate with third parties');

INSERT INTO servicetable (id, name)
VALUES (12, 'Contact Us');

INSERT INTO servicetable (id, name)
VALUES (13, 'Payments');

INSERT INTO servicetable (id, name)
VALUES (14, 'Oauth');

INSERT INTO servicetable (id, name, description)
VALUES (15, 'Tracing', 'Improve visibility of your distributed services');

INSERT INTO servicetable (id, name)
VALUES (16, 'Authenticate');

INSERT INTO servicetable (id, name, description)
VALUES (17, 'Payments', 'Storefront for your business');