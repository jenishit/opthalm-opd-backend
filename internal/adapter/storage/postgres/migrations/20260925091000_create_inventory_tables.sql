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
