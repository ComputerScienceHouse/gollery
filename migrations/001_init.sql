CREATE TABLE directory (
    did                 serial PRIMARY KEY,
    name                varchar(32),
    creator             varchar(32),
    create_date         date
);

CREATE TABLE subdirectory (
    parent_did          int,
    sub_did             int,
    PRIMARY KEY (parent_did, sub_did),
    FOREIGN KEY (parent_did) REFERENCES directory(did),
    FOREIGN KEY (sub_did) REFERENCES directory(did)
);


CREATE TABLE groups (
    gid                 serial PRIMARY KEY,
    group_name          varchar(30)
);

CREATE TABLE file (
    fid                 serial,
    name                varchar(128),
    creator             char(32),
    upload_time         date,
    s3_key              varchar(100),
    did                 int,
    PRIMARY KEY (fid),
    FOREIGN KEY (did) REFERENCES directory(did)
);

CREATE TABLE file_groups (
    fid                 int REFERENCES file(fid),
    gid                 int REFERENCES groups(gid),
    PRIMARY KEY (fid, gid)
);