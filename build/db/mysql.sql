create table msnger.Operation
(
    revisionId bigint unsigned NOT NULL,
    type       int       NOT NULL,
    param1     varchar(50),
    param2     varchar(50),
    param3     varchar(50),
    messageId  varchar(21),

    createdAt  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    primary key (revisionId)
);

create table msnger.Message
(
    id          varchar(21) NOT NULL,
    `to`        varchar(21) NOT NULL,
    `from`      varchar(21) NOT NULL,
    contentType int         NOT NULL,
    text        varchar(1000),
    metadata    json        NOT NULL DEFAULT (JSON_OBJECT()),

    createdAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    primary key (id)
);

create table msnger.User
(
    id          varchar(21) NOT NULL,
    displayName varchar(30) NOT NULL,
    statusText  text,
    pictureUrl  text,

    createdAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    primary key (id)
);

create table msnger.Account
(
    id varchar(21) NOT NULL,
    email varchar(300) NOT NULL,
    password text NOT NULL,
    isAdmin bit NOT NULL DEFAULT 0,

    createdAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt   timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    primary key(id),
    UNIQUE(email)
);

create table msnger.OpRelation
(
    id           int auto_increment NOT NULL,
    revisionId   bigint unsigned NOT NULL,
    targetUserId varchar(21) NOT NULL,

    primary key (id),
    UNIQUE (revisionId, targetUserId)
);

create table msnger.LastRevision
(
    id varchar(21) NOT NULL,
    lastRevisionId bigint unsigned NOT NULL DEFAULT 0,

    primary key (id)
);