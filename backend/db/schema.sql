CREATE TABLE word (
    id int NOT NULL PRIMARY KEY,
    name TEXT,
    name_stressed TEXT,
    name_broken TEXT,
    type_id int
);

CREATE TABLE word_type (
    id int NOT NULL PRIMARY KEY,
    name TEXT,
    speech_part TEXT
);

CREATE TABLE derivative_form (
    id int NOT NULL PRIMARY KEY,
    name TEXT,
    name_stressed TEXT,
    name_broken TEXT,
    name_condensed TEXT,
    description TEXT,
    is_infinitive int,
    base_word_id int
);

CREATE TABLE incorrect_form (
    id int,
    name TEXT,
    correct_word_id int
);

create TABLE word_translation (
    id int NOT NULL PRIMARY KEY,
    word_id int,
    lang TEXT,
    content TEXT
)