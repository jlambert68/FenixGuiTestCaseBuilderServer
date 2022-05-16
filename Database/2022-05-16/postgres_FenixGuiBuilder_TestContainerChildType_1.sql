create table "TestContainerChildType"
(
    "TestInstructionContainerChildType"        integer not null
        constraint testcontainerchildtype_pk
            unique,
    "TestInstructionContainerChildDescription" varchar not null
);

comment on table "TestContainerChildType" is 'Holds the different child-types that a TestInstructionContainer can have';

alter table "TestContainerChildType"
    owner to postgres;

INSERT INTO "FenixGuiBuilder"."TestContainerChildType" ("TestInstructionContainerChildType", "TestInstructionContainerChildDescription") VALUES (1, 'TestInstruction');
INSERT INTO "FenixGuiBuilder"."TestContainerChildType" ("TestInstructionContainerChildType", "TestInstructionContainerChildDescription") VALUES (2, 'TestContainer');
