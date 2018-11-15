-- +migrate Up
create type devices_status as enum ('enabled', 'disabled');

CREATE TABLE devices (
  id          bigserial constraint devices_pkey primary key not null,
  name        VarChar(255)                                  NOT NULL,
  description VarChar(255)                                  default '',
  device_id   smallint                                      NULL,
  node_id     BIGINT CONSTRAINT devices_2_nodes_fk REFERENCES nodes (id) on update cascade on delete cascade null,
  baud        smallint                                      default 0,
  tty         VarChar(255)                                  NULL,
  stop_bite   smallint                                      default 0,
  timeout     smallint                                      default 0,
  address     smallint                                      NULL,
  status      devices_status                                NOT NULL DEFAULT 'enabled',
  sleep       smallint                                      NOT NULL DEFAULT 0,
  created_at  timestamp with time zone                      not null,
  updated_at  timestamp with time zone                      not null
);

CREATE UNIQUE INDEX name_device_address_2_devices_unq ON devices (name, device_id, address);
CREATE UNIQUE INDEX name_address_2_devices_unq ON devices (node_id, address);
CREATE UNIQUE INDEX device_address_2_devices_unq ON devices (device_id, address);

CREATE TABLE device_actions (
  id          bigserial constraint device_actions_pkey primary key                                                      NOT NULL,
  device_id   BIGINT CONSTRAINT device_actions_2_devices_fk REFERENCES devices (id) on update cascade on delete cascade,
  name        VarChar(255)                                                                                              NOT NULL,
  description VarChar(255)                                                                                              NULL,
  script_id   BIGINT CONSTRAINT device_actions_2_scripts_fk REFERENCES scripts (id) on update cascade on delete cascade NULL,
  created_at  timestamp with time zone                                                                                  NOT NULL,
  updated_at  timestamp with time zone                                                                                  NOT NULL
);

CREATE UNIQUE INDEX device_name_2_device_actions_unq ON device_actions (device_id, name);

CREATE TABLE device_states (
  id          bigserial                not null constraint device_states_pkey primary key,
  system_name VarChar(255)             NOT NULL,
  description Text                     NOT NULL,
  device_id   BIGINT CONSTRAINT device_states_2_devices_fk REFERENCES devices (id) on update cascade on delete cascade,
  created_at  timestamp with time zone not null,
  updated_at  timestamp with time zone not null
);

CREATE UNIQUE INDEX device_nsystem_name_2_device_states_unq ON device_states (device_id, system_name);

-- +migrate Down
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS device_actions CASCADE;
DROP TABLE IF EXISTS device_states CASCADE;
drop type devices_status;