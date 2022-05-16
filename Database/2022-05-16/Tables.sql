
-- auto-generated definition
create schema "FenixGuiBuilder";

alter schema "FenixGuiBuilder" owner to postgres;





create table if not exists "FenixGuiBuilder"."TestInstructions"
(
    "DomainUuid"                   uuid      not null
    constraint testinstructions_domains_domain_uuid_fk
    references domains,
    "DomainName"                   varchar   not null,
    "TestInstructionUuid"          uuid      not null
    constraint testinstructions_pk
    primary key,
    "TestInstructionName"          varchar   not null,
    "TestInstructionTypeUuid"      uuid      not null,
    "TestInstructionTypeName"      varchar   not null,
    "TestInstructionDescription"   varchar   not null,
    "TestInstructionMouseOverText" varchar   not null,
    "Deprecated"                   boolean   not null,
    "Enabled"                      boolean   not null,
    "MajorVersionNumber"           integer   not null,
    "MinorVersionNumber"           integer   not null,
    "UpdatedTimeStamp"             timestamp not null
);

comment on table "FenixGuiBuilder"."TestInstructions" is 'Holds all TestInstructions';

alter table "FenixGuiBuilder"."TestInstructions"
    owner to postgres;

create table if not exists "FenixGuiBuilder"."TestInstructionContainers"
(
    "DomainUuid"                            uuid      not null,
    "DomainName"                            varchar   not null,
    "TestInstructionContainerUuid"          uuid      not null
    constraint testinstructioncontainers_pk
    primary key,
    "TestInstructionContainerName"          varchar   not null,
    "TestInstructionContainerTypeUuid"      uuid      not null,
    "TestInstructionContainerTypeName"      varchar   not null,
    "TestInstructionContainerDescription"   varchar   not null,
    "TestInstructionContainerMouseOverText" varchar   not null,
    "Deprecated"                            boolean   not null,
    "Enabled"                               boolean   not null,
    "MajorVersionNumber"                    integer   not null,
    "MinorVersionNumber"                    integer   not null,
    "UpdatedTimeStamp"                      timestamp not null,
    "ChildrenIsParallelProcessed"           boolean   not null
);

alter table "FenixGuiBuilder"."TestInstructionContainers"
    owner to postgres;

create table if not exists "FenixGuiBuilder"."TestContainerChildType"
(
    "TestInstructionContainerChildType"        integer not null
    constraint testcontainerchildtype_pk
    unique,
    "TestInstructionContainerChildDescription" varchar not null
);

comment on table "FenixGuiBuilder"."TestContainerChildType" is 'Holds the different child-types that a TestInstructionContainer can have';

alter table "FenixGuiBuilder"."TestContainerChildType"
    owner to postgres;

create table if not exists "FenixGuiBuilder"."TestInstructionContainerChildren"
(
    "TestIntructionContainer" uuid    not null
    constraint testinstructioncontainerchildren_testinstructioncontainercontai
    references "FenixGuiBuilder"."TestInstructionContainers",
    "TestContainerChildUuid"  uuid    not null,
    "TestContainerChildType"  integer not null
    constraint testinstructioncontainerchildren_testcontainerchildtype_testins
    references "FenixGuiBuilder"."TestContainerChildType" ("TestInstructionContainerChildType")
    );

comment on table "FenixGuiBuilder"."TestInstructionContainerChildren" is 'Holds the releation between a TestIntructionContainer and it''s children';

alter table "FenixGuiBuilder"."TestInstructionContainerChildren"
    owner to postgres;

create table if not exists "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers"
(
    "UserId"     varchar   not null,
    "PinnedUuid" uuid      not null,
    "PinnedName" varchar   not null,
    "PinnedType" integer   not null
    constraint pinnedtestinstructionsandprecreatedtestinstructioncontainers_te
    references "FenixGuiBuilder"."TestContainerChildType" ("TestInstructionContainerChildType"),
    "TimeStamp"  timestamp not null,
    constraint pinnedtestinstructionsandprecreatedtestinstructioncontainers_pk
    primary key ("UserId", "PinnedUuid")
    );

comment on table "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" is 'Holds all users pinned TestInstructions ans pre-created TestInstructionsContainers';

alter table "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers"
    owner to postgres;

create unique index if not exists pinnedtestinstructionsandprecreatedtestinstructioncontainers_pi
    on "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" ("PinnedUuid");

