create table "TestInstructionContainers"
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

alter table "TestInstructionContainers"
    owner to postgres;

INSERT INTO "FenixGuiBuilder"."TestInstructionContainers" ("DomainUuid", "DomainName", "TestInstructionContainerUuid", "TestInstructionContainerName", "TestInstructionContainerTypeUuid", "TestInstructionContainerTypeName", "TestInstructionContainerDescription", "TestInstructionContainerMouseOverText", "Deprecated", "Enabled", "MajorVersionNumber", "MinorVersionNumber", "UpdatedTimeStamp", "ChildrenIsParallelProcessed") VALUES ('e81b9734-5dce-43c9-8d77-3368940cf126', 'Fenix', 'b107bdd9-4152-4020-b3f0-fc750b45885e', 'Emtpy parallelled processed TestInstructionsContainer', 'b107bdd9-4152-4020-b3f0-fc750b45885e', 'Base containers', 'Children of this container is processed in parallel', 'Children of this container is processed in parallel', false, true, 0, 1, '2022-05-02 10:08:28.000000', true);
INSERT INTO "FenixGuiBuilder"."TestInstructionContainers" ("DomainUuid", "DomainName", "TestInstructionContainerUuid", "TestInstructionContainerName", "TestInstructionContainerTypeUuid", "TestInstructionContainerTypeName", "TestInstructionContainerDescription", "TestInstructionContainerMouseOverText", "Deprecated", "Enabled", "MajorVersionNumber", "MinorVersionNumber", "UpdatedTimeStamp", "ChildrenIsParallelProcessed") VALUES ('e81b9734-5dce-43c9-8d77-3368940cf126', 'Fenix', 'e81b9734-5dce-43c9-8d77-3368940cf126', 'Emtpy serial processed TestInstructionsContainer', 'b107bdd9-4152-4020-b3f0-fc750b45885e', 'Base containers', 'Children of this container is processed in serial', 'Children of this container is processed in serial', false, true, 0, 1, '2022-05-02 10:10:07.000000', false);
