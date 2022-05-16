create table "TestInstructionContainerChildren"
(
    "TestIntructionContainer" uuid    not null
        constraint testinstructioncontainerchildren_testinstructioncontainercontai
            references "TestInstructionContainers",
    "TestContainerChildUuid"  uuid    not null,
    "TestContainerChildType"  integer not null
        constraint testinstructioncontainerchildren_testcontainerchildtype_testins
            references "TestContainerChildType" ("TestInstructionContainerChildType")
);

comment on table "TestInstructionContainerChildren" is 'Holds the releation between a TestIntructionContainer and it''s children';

alter table "TestInstructionContainerChildren"
    owner to postgres;

