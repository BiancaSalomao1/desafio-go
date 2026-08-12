CREATE TABLE order_items (

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL,

    product_id UUID NOT NULL,

    product_name VARCHAR(150) NOT NULL,

    product_price NUMERIC(10,2) NOT NULL,

    quantity INTEGER NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_order
        FOREIGN KEY(order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
);