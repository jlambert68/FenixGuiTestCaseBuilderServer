create table "PinnedTestInstructionsAndPreCreatedTestInstructionContainers"
(
    "UserId"     varchar   not null,
    "PinnedUuid" uuid      not null,
    "PinnedName" varchar   not null,
    "PinnedType" integer   not null
        constraint pinnedtestinstructionsandprecreatedtestinstructioncontainers_te
            references "TestContainerChildType" ("TestInstructionContainerChildType"),
    "TimeStamp"  timestamp not null,
    constraint pinnedtestinstructionsandprecreatedtestinstructioncontainers_pk
        primary key ("UserId", "PinnedUuid")
);

comment on table "PinnedTestInstructionsAndPreCreatedTestInstructionContainers" is 'Holds all users pinned TestInstructions ans pre-created TestInstructionsContainers';

alter table "PinnedTestInstructionsAndPreCreatedTestInstructionContainers"
    owner to postgres;

create unique index pinnedtestinstructionsandprecreatedtestinstructioncontainers_pi
    on "PinnedTestInstructionsAndPreCreatedTestInstructionContainers" ("PinnedUuid");

INSERT INTO "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" ("UserId", "PinnedUuid", "PinnedName", "PinnedType", "TimeStamp") VALUES ('s41797', '2f130d7e-f8aa-466f-b29d-0fb63608c1a6', 'TestInstructionName 1', 1, '2022-05-11 11:05:35.240448');
INSERT INTO "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" ("UserId", "PinnedUuid", "PinnedName", "PinnedType", "TimeStamp") VALUES ('s41797', 'b107bdd9-4152-4020-b3f0-fc750b45885e', 'TestInstructionContainerName 1', 2, '2022-05-11 11:05:35.240448');
INSERT INTO "FenixGuiBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" ("UserId", "PinnedUuid", "PinnedName", "PinnedType", "TimeStamp") VALUES ('s41797', 'e81b9734-5dce-43c9-8d77-3368940cf126', 'TestInstructionContainerName', 2, '2022-05-11 11:05:35.240448');
