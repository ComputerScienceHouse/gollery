CREATE TABLE directory (
    did                 char(32),
    name                varchar(32)
);

CREATE TABLE subdirectory (
    parent_did          char(32),
    sub_did             char(32),
    PRIMARY KEY (parent_did, sub_did),
    FOREIGN KEY (parent_did) REFERENCES directory(did),
    FOREIGN KEY (sub_did) REFERENCES directory(did)
);


CREATE TABLE groups (
    gid                 char(32),
    group_name          varchar(30)
);

CREATE TABLE file_groups (
    fid                 int REFERENCES file(fid),
    gid                 int REFERENCES groups(gid),
    PRIMARY KEY (fid, gid)
);

CREATE TABLE file (
    fid                 serial,
    uuid                char(32),
    upload_time         date,  
    s3_key              varchar(100),
    did                 char(32),
    PRIMARY KEY (fid),
    FOREIGN KEY (did) REFERENCES directory(did)
);