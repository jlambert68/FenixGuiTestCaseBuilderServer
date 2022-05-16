create table "TestInstructions"
(
    "DomainUuid"                   uuid      not null
        constraint testinstructions_domains_domain_uuid_fk
            references public.domains,
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

comment on table "TestInstructions" is 'Holds all TestInstructions';

alter table "TestInstructions"
    owner to postgres;

INSERT INTO "FenixGuiBuilder"."TestInstructions" ("DomainUuid", "DomainName", "TestInstructionUuid", "TestInstructionName", "TestInstructionTypeUuid", "TestInstructionTypeName", "TestInstructionDescription", "TestInstructionMouseOverText", "Deprecated", "Enabled", "MajorVersionNumber", "MinorVersionNumber", "UpdatedTimeStamp") VALUES ('78a97c41-a098-4122-88d2-01ed4b6c4844', 'Custody Arrangement', '2f130d7e-f8aa-466f-b29d-0fb63608c1a6', 'Just the name', '513dd8fb-a0bb-4738-9a0b-b7eaf7bb8adb', 'The type of the TestInstruction', 'En vanlig typ', 'This will be shown when hovering above this TestInstruction', false, true, 0, 1, '2022-04-29 15:42:15.000000');
