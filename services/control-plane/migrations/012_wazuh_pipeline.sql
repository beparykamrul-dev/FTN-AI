INSERT INTO permissions(key, description) VALUES
 ('security.alert.ingest','Ingest normalized Wazuh security alerts from FTN-owned infrastructure')
ON CONFLICT (key) DO NOTHING;
