-- +goose Up
CREATE TABLE inventory_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category          VARCHAR(20) NOT NULL
                          CHECK (category IN ('frame', 'lens', 'contact_lens', 'sunglasses', 'coating', 'accessory', 'other')),
    sku               VARCHAR(50) NOT NULL UNIQUE,
    name              VARCHAR(150) NOT NULL,
    brand             VARCHAR(100),
    model             VARCHAR(100),
    color             VARCHAR(50),
    size              VARCHAR(50),
    cost_price        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    selling_price     NUMERIC(10, 2) NOT NULL DEFAULT 0,
    quantity_on_hand  INTEGER NOT NULL DEFAULT 0,
    reorder_threshold INTEGER NOT NULL DEFAULT 5,
    unit              VARCHAR(20) NOT NULL DEFAULT 'pcs',
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at        TIMESTAMP DEFAULT NULL,
    created_by        UUID NOT NULL REFERENCES users(id),
    updated_by        UUID NOT NULL REFERENCES users(id),
    created_at        TIMESTAMP NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE inventory_stock_movements (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id),
    movement_type     VARCHAR(20) NOT NULL
                          CHECK (movement_type IN ('purchase_in', 'sale_out', 'adjustment_in', 'adjustment_out', 'return_in')),
    quantity          INTEGER NOT NULL CHECK (quantity > 0),
    reference_type    VARCHAR(20),
    reference_id      UUID,
    notes             VARCHAR(255),
    created_by        UUID NOT NULL REFERENCES users(id),
    created_at        TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_movements_item_id ON inventory_stock_movements (inventory_item_id);

CREATE TABLE vendors (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(150) NOT NULL,
    contact_person  VARCHAR(100),
    phone           VARCHAR(20),
    email           VARCHAR(150),
    address         VARCHAR(255),
    deleted_at      TIMESTAMP DEFAULT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE stock_purchases (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id     UUID NOT NULL REFERENCES vendors(id),
    purchase_date DATE NOT NULL DEFAULT CURRENT_DATE,
    invoice_ref_no VARCHAR(100),
    total_amount  NUMERIC(10, 2) NOT NULL DEFAULT 0,
    paid_amount   NUMERIC(10, 2) NOT NULL DEFAULT 0,
    due_amount    NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_by    UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMP NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE stock_purchase_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_id       UUID NOT NULL REFERENCES stock_purchases(id),
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id),
    quantity          INTEGER NOT NULL CHECK (quantity > 0),
    unit_cost         NUMERIC(10, 2) NOT NULL DEFAULT 0,
    line_total        NUMERIC(10, 2) NOT NULL DEFAULT 0
);

CREATE INDEX idx_stock_purchase_items_purchase_id ON stock_purchase_items (purchase_id);

ALTER TABLE invoice_items
    ADD CONSTRAINT fk_invoice_items_inventory_item
    FOREIGN KEY (inventory_item_id) REFERENCES inventory_items(id);

-- +goose Down
ALTER TABLE invoice_items DROP CONSTRAINT fk_invoice_items_inventory_item;
DROP TABLE stock_purchase_items;
DROP TABLE stock_purchases;
DROP TABLE vendors;
DROP TABLE inventory_stock_movements;
DROP TABLE inventory_items;
