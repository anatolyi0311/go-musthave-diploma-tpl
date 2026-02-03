CREATE TABLE users(
    id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    create_time timestamp without time zone DEFAULT '2026-01-10 19:28:13.459931'::timestamp without time zone,
    login varchar(255),
    password varchar(255),
    PRIMARY KEY(id)
);

CREATE TABLE user_balance(
    id SERIAL NOT NULL,
    "current" double precision,
    withdrawn double precision,
    user_id integer,
    updated_at timestamp without time zone NOT NULL,
    uploaded_at timestamp without time zone NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT user_balance_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);

CREATE TABLE orders(
    number varchar(50) NOT NULL,
    accrual double precision,
    user_id integer,
    status varchar(20) NOT NULL,
    uploaded_at timestamp without time zone NOT NULL,
    processed_at timestamp without time zone,
    PRIMARY KEY(number),
    CONSTRAINT orders_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);

CREATE TABLE balance_operations(
    id SERIAL NOT NULL,
    user_id integer,
    "type" varchar(20) NOT NULL,
    amount double precision NOT NULL,
    order_number varchar(50),
    processed_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT balance_operations_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);