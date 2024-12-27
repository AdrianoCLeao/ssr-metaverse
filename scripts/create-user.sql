CREATE TABLE users (
    id SERIAL PRIMARY KEY,              
    username VARCHAR(50) NOT NULL,     
    password TEXT NOT NULL,             
    created_at TIMESTAMP DEFAULT NOW()  
);
