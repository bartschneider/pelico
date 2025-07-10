-- PostgreSQL initialization script for Pelico

-- Create database if it doesn't exist (handled by POSTGRES_DB env var)

-- Create user if it doesn't exist (handled by POSTGRES_USER env var)

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Initial data for common platforms
INSERT INTO platforms (name, manufacturer, release_year, created_at, updated_at) VALUES
    ('Nintendo Entertainment System', 'Nintendo', 1985, NOW(), NOW()),
    ('Super Nintendo Entertainment System', 'Nintendo', 1991, NOW(), NOW()),
    ('Nintendo 64', 'Nintendo', 1996, NOW(), NOW()),
    ('Nintendo GameCube', 'Nintendo', 2001, NOW(), NOW()),
    ('Nintendo Wii', 'Nintendo', 2006, NOW(), NOW()),
    ('Game Boy', 'Nintendo', 1989, NOW(), NOW()),
    ('Game Boy Color', 'Nintendo', 1998, NOW(), NOW()),
    ('Game Boy Advance', 'Nintendo', 2001, NOW(), NOW()),
    ('Nintendo DS', 'Nintendo', 2004, NOW(), NOW()),
    ('Nintendo 3DS', 'Nintendo', 2011, NOW(), NOW()),
    ('Sega Master System', 'Sega', 1986, NOW(), NOW()),
    ('Sega Genesis', 'Sega', 1988, NOW(), NOW()),
    ('Sega Saturn', 'Sega', 1994, NOW(), NOW()),
    ('Sega Dreamcast', 'Sega', 1998, NOW(), NOW()),
    ('Sony PlayStation', 'Sony', 1994, NOW(), NOW()),
    ('Sony PlayStation 2', 'Sony', 2000, NOW(), NOW()),
    ('Sony PlayStation 3', 'Sony', 2006, NOW(), NOW()),
    ('Sony PlayStation 4', 'Sony', 2013, NOW(), NOW()),
    ('Sony PlayStation 5', 'Sony', 2020, NOW(), NOW()),
    ('Sony PlayStation Portable', 'Sony', 2004, NOW(), NOW()),
    ('Sony PlayStation Vita', 'Sony', 2011, NOW(), NOW()),
    ('Microsoft Xbox', 'Microsoft', 2001, NOW(), NOW()),
    ('Microsoft Xbox 360', 'Microsoft', 2005, NOW(), NOW()),
    ('Microsoft Xbox One', 'Microsoft', 2013, NOW(), NOW()),
    ('Microsoft Xbox Series X/S', 'Microsoft', 2020, NOW(), NOW()),
    ('Atari 2600', 'Atari', 1977, NOW(), NOW()),
    ('Atari 5200', 'Atari', 1982, NOW(), NOW()),
    ('Atari 7800', 'Atari', 1986, NOW(), NOW()),
    ('PC', 'Various', 1981, NOW(), NOW()),
    ('Arcade', 'Various', 1970, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;