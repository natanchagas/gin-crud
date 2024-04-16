USE real_states;

CREATE TABLE real_states (
    real_state_id INT AUTO_INCREMENT PRIMARY KEY,
    real_state_registration INT(100) UNIQUE NOT NULL,
    real_state_address TEXT NOT NULL,
    real_state_size DECIMAL(10,2) NOT NULL,
    real_state_price DECIMAL(15,2) NOT NULL,
    real_state_state VARCHAR(2) NOT NULL
);

