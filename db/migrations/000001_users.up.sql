-- Create enum
CREATE TYPE enum_height_units as ENUM ('CM', 'INCH');
CREATE TYPE enum_weight_units as ENUM ('KG', 'LBS');
CREATE TYPE enum_preferences as ENUM ('CARDIO', 'WEIGHT');

-- Create table users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    user_image_uri VARCHAR(255),
    weight INTEGER,
    height INTEGER,
    weight_unit enum_weight_units,
    height_unit enum_height_units,
    preference enum_preferences
);
