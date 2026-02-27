-- +goose Up
-- +goose StatementBegin

CREATE TABLE books(
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    publication_year INTEGER NOT NULL CHECK(publication_year > 0),
    pages INTEGER NOT NULL CHECK(pages > 0),
    genre TEXT NOT NULL
);

-- Фэнтези (#1)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'J.R.R. Tolkien',
    'The Hobbit',
    1937,
    310,
    'Fantasy'
);

-- Фэнтези (#2)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'J.R.R. Tolkien',
    'The Lord of the Rings: The Fellowship of the Ring',
    1954,
    423,
    'Fantasy'
);

-- Фэнтези (#3)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'George R.R. Martin',
    'A Game of Thrones',
    1996,
    694,
    'Fantasy'
);

-- Фэнтези (#4)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Neil Gaiman',
    'American Gods',
    2001,
    634,
    'Fantasy'
);

-- Фэнтези (#5)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Terry Pratchett',
    'The Colour of Magic',
    1983,
    208,
    'Fantasy'
);

-- Научная фантастика (#1)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Isaac Asimov',
    'Foundation',
    1951,
    256,
    'Science Fiction'
);

-- Научная фантастика (#2)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Frank Herbert',
    'Dune',
    1965,
    544,
    'Science Fiction'
);

-- Научная фантастика (#3)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Ray Bradbury',
    'Fahrenheit 451',
    1953,
    240,
    'Science Fiction'
);

-- Научная фантастика (#4)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Douglas Adams',
    'The Hitchhikers Guide to the Galaxy',
    1979,
    224,
    'Science Fiction'
);

-- Научная фантастика (#5)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Philip K. Dick',
    'Do Androids Dream of Electric Sheep?',
    1968,
    288,
    'Science Fiction'
);

-- Классика (#1)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Leo Tolstoy',
    'War and Peace',
    1869,
    1225,
    'Classic'
);

-- Классика (#2)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Fyodor Dostoevsky',
    'Crime and Punishment',
    1866,
    544,
    'Classic'
);

-- Классика (#3)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Charles Dickens',
    'Great Expectations',
    1861,
    544,
    'Classic'
);

-- Классика (#4)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Jane Austen',
    'Pride and Prejudice',
    1813,
    279,
    'Classic'
);

-- Классика (#5)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Franz Kafka',
    'The Metamorphosis',
    1915,
    128,
    'Classic'
);

-- Детектив (#1)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Agatha Christie',
    'Murder on the Orient Express',
    1934,
    256,
    'Detective'
);

-- Детектив (#2)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Sir Arthur Conan Doyle',
    'The Hound of the Baskervilles',
    1902,
    256,
    'Detective'
);

-- Детектив (#3)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Raymond Chandler',
    'The Big Sleep',
    1939,
    256,
    'Detective'
);

-- Детектив (#4)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'Dashiell Hammett',
    'The Maltese Falcon',
    1930,
    256,
    'Detective'
);

-- Детектив (#5)
INSERT INTO books (
    author, 
    title, 
    publication_year,
    pages,
    genre
) 
VALUES (
    'James Patterson',
    'Along Came a Spider',
    1993,
    352,
    'Detective'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
